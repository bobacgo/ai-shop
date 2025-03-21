package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	v1 "github.com/bobacgo/ai-shop/api/pb/user/v1"
	"github.com/bobacgo/ai-shop/api/pb/user/v1/errs"
	"github.com/bobacgo/ai-shop/user/internal/config"
	"github.com/bobacgo/ai-shop/user/internal/repo"
	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/kit/app/security"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type AuthService struct {
	v1.UnimplementedAuthServiceServer
	b64c             *base64Captcha.Captcha
	passwordStrength *security.PasswordValidator // 密码强度校验
	passwdVerifier   *security.PasswdVerifier    // 登录密码验证器
	jwtHelper        *security.JWToken
	repo             *repo.All
}

func NewAuthService(rdb redis.UniversalClient, repo *repo.All) *AuthService {
	return &AuthService{
		b64c: base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, repo.Captcha),
		// 密码强度: 大小写+数字、 len > 6
		passwordStrength: security.NewPasswordValidator(6, true, true, true, false),
		passwdVerifier:   security.NewPasswdVerifier(rdb, config.Cfg().Service.ErrAttemptLimit),
		jwtHelper:        security.NewJWT(&config.Cfg().Security.Jwt, rdb),
		repo:             repo,
	}
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {

	// 要求
	// 1.校验验证码, 并清空验证码
	// 2.对密码字段从密文转换成明文
	// 3.是否禁用（包括已注销的）
	// 4.校验密码
	// 5.登录次数（防止被暴力破解），登录次数超过5次直接禁用
	// 6.重置token（同时只能一个地方登录），如果有
	// 7.更新登录日志

	// 1.支持多平台， 不同平台有不同的 token
	// 2.每一个平台只能登录同时在线一个

	// 1. 校验验证码，并清空验证码
	if !s.repo.Captcha.Verify(req.VerificationKey, req.VerificationCode, true) {
		return nil, errs.Status(ctx, errs.Err_CaptchaValidFailed)
	}

	// 2. 对密码字段从密文转换成明文
	ciphertext := security.Ciphertext(config.Cfg().Security.Ciphertext.CipherKey)
	if err := ciphertext.Decrypt(req.Password); err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt",
			"secret", security.Ciphertext(config.Cfg().Security.Ciphertext.CipherKey), // 脱敏打印
			"password", req.Password,
			"err", err,
		)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}
	req.Password = string(ciphertext)

	// 3. 查询用户信息并检查状态
	user, err := s.repo.User.FindOneUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// TODO 记录登录失败次数
			return nil, errs.Status(ctx, errs.Err_UsernameOrPasswordErr)
		}
		slog.ErrorContext(ctx, "s.repo.UserRepo.FindOneUserByUsername", "username", req.Username, "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}
	if user.Status == 0 {
		return nil, errs.Status(ctx, errs.Err_UserBanned)
	} else if user.Status == 2 {
		return nil, errs.Status(ctx, errs.Err_UserDeleting)
	}

	// 3. 校验密码
	// 4. TODO 登录次数（防止被暴力破解），登录次数超过5次直接禁用
	if !s.passwdVerifier.BcryptVerifyWithCount(ctx, string(user.Password), req.Password) {
		return nil, errs.Status(ctx, errs.Err_UsernameOrPasswordErr)
	}

	// 5. 生成新token
	accessToken, refreshToken, err := s.jwtHelper.Generate(ctx, &security.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    []string{strconv.Itoa(int(user.Role))},
	})
	if err != nil {
		return nil, errs.Status(ctx, errs.Err_LoginFailed)
	}

	// 6. 重置token（同时只能一个地方登录），如果有

	// 7. 更新登录日志
	loginLog := &model.UserLoginSuccessLog{
		UserID:            user.ID,
		Username:          user.Username,
		LoginTime:         time.Now(),
		ClientIP:          "",
		UserAgent:         nil,
		LoginChannel:      0,
		GeoLocation:       nil,
		DeviceFingerprint: nil,
		RiskLevel:         0,
	}
	if err := s.repo.UserLoginSuccessLog.InsertUserLoginSuccessLog(ctx, loginLog); err != nil {
		slog.ErrorContext(ctx, "s.repo.UserLoginSuccessLog.InsertUserLoginSuccessLog", "err", err, "loginLog", loginLog)
	}

	return &v1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    config.Cfg().Basic.Security.Jwt.AccessTokenExpired.TimeDuration().Microseconds(),
	}, nil
}

func (s *AuthService) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.LoginResponse, error) {

	// 要求
	// 1. 校验验证码 (保留验证码登录可以继续使用)
	// 2. 对密码字段从密文转换成明文
	// 3. 码强度校验
	// 4. 检验用户是否被占用（用户名，手机号， 邮箱号）（包括被禁用的和已注销的）
	// 5. 密码字段加密
	// 6. 保存用户信息
	// 7. 登录并返回token

	// 1. 校验验证码 (保留验证码登录可以继续使用)
	if !s.b64c.Verify(request.VerificationKey, request.VerificationCode, false) {
		return nil, errs.Status(ctx, errs.Err_CaptchaValidFailed)
	}

	// 2. 对密码字段从密文转换成明文
	ciphertext := security.Ciphertext(config.Cfg().Security.Ciphertext.CipherKey)
	if err := ciphertext.Decrypt(request.Password); err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt",
			"secret", security.Ciphertext(config.Cfg().Security.Ciphertext.CipherKey), // 脱敏打印
			"password", request.Password,
			"err", err,
		)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 3. 密码强度校验
	if level, err := s.passwordStrength.Validate(string(ciphertext)); err != nil {
		st, _ := status.New(codes.Code(errs.Err_InvalidPasswordFormat),
			errs.GetErrorMessage(ctx, errs.Err_InvalidPasswordFormat)).
			WithDetails(&v1.PasswordStrengthError{Level: int32(level)})
		return nil, st.Err()
	}

	// 4. 检验用户名是否被占用
	if _, err := s.repo.User.FindOneUserByUsername(ctx, request.Username); err == nil {
		return nil, errs.Status(ctx, errs.Err_UserAlreadyExists)
	}

	// 5. 密码字段加密
	hashedPassword := ciphertext.BcryptHash()

	// 6. 保存用户信息
	user := &model.User{
		ID:       uuid.New().String(),
		Username: request.Username,
		Email:    &request.Email,
		Phone:    &request.Phone,
		Password: []byte(hashedPassword),
		Role:     int32(v1.Role_customer),     // 默认为消费者角色
		Status:   int32(v1.UserStatus_active), // 默认为正常状态
	}

	if err := s.repo.User.InsertUser(ctx, user); err != nil {
		slog.ErrorContext(ctx, "s.repo.UserRepo.InsertUser", "err", err, "user", user)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}

	// 7. 登录并返回token
	return s.Login(ctx, &v1.LoginRequest{
		Username:         request.Username,
		Password:         request.Password,
		VerificationCode: request.VerificationCode,
		VerificationKey:  request.VerificationKey,
	})
}

func (s *AuthService) ResetPassword(ctx context.Context, request *v1.ResetPasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errs.Status(ctx, errs.Err_UserNotFound)
}

// 生成验证码
func (s *AuthService) SendVerificationCode(ctx context.Context, _ *emptypb.Empty) (*v1.SendVerificationCodeResponse, error) {
	captchaKey, imgBase64, answer, err := s.b64c.Generate()
	if err != nil {
		slog.ErrorContext(ctx, "s.b64c.Generate()", "err", err)
		return nil, errs.Status(ctx, errs.Err_CaptchaLoad)
	}
	slog.DebugContext(ctx, "captcha", "captchaKey", captchaKey, "answer", answer)
	return &v1.SendVerificationCodeResponse{
		VerificationKey:   captchaKey,
		VerificationImage: imgBase64,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, request *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	return nil, errs.Status(ctx, errs.Err_UserNotFound)
}

func (s *AuthService) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errs.Status(ctx, errs.Err_UserNotFound)
}

func (s *AuthService) Deregister(ctx context.Context, request *v1.DeregisterRequest) (*emptypb.Empty, error) {

	// 要求
	// 1. 对密码字段从密文转换成明文
	// 2. 校验用户密码
	// 3. 更改用户状态为注销, 并添加注销请求记录
	// 4. 注销用户的所有token

	ciphertext := security.Ciphertext(config.Cfg().Security.Ciphertext.CipherKey)
	if err := ciphertext.Decrypt(request.Password); err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt",
			"secret", security.Ciphertext(config.Cfg().Security.Ciphertext.CipherKey), // 脱敏打印
			"password", request.Password,
			"err", err,
		)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 验证用户密码
	user, err := s.repo.User.FindOneUserById(ctx, request.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Status(ctx, errs.Err_UserNotFound)
		}
		slog.ErrorContext(ctx, "s.repo.UserRepo.FindOneUserById", "userId", request.UserId, "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}

	if !s.passwdVerifier.BcryptVerify(string(user.Password), string(ciphertext)) {
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 更改用户状态为注销
	// 并添加注销请求记录
	if err := s.repo.User.DeletedUser(ctx, request.UserId); err != nil {
		slog.ErrorContext(ctx, "s.repo.UserRepo.DeleteUser", "userId", request.UserId, "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}

	// 注销用户的所有token
	if err = s.jwtHelper.RemoveToken(ctx, user.Username); err != nil {
		slog.ErrorContext(ctx, "s.jwtHelper.Deregister", "username", user.Username, "err", err)
		return nil, nil
	}
	return &emptypb.Empty{}, nil
}
