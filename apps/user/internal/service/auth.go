package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	v1 "github.com/bobacgo/ai-shop/api/pb/user/v1"
	"github.com/bobacgo/ai-shop/api/pb/user/v1/errs"
	"github.com/bobacgo/ai-shop/user/internal/config"
	"github.com/bobacgo/ai-shop/user/internal/repo"
	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/kit/app/security"
	"github.com/bobacgo/kit/pkg/ucrypto"
	"github.com/golang-jwt/jwt/v5"
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

func NewAuthService(rdb redis.UniversalClient, r *repo.All) *AuthService {
	return &AuthService{
		b64c: base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, r.Captcha),
		// 密码强度: 大小写+数字、 len > 6
		passwordStrength: security.NewPasswordValidator(6, true, true, true, false),
		passwdVerifier:   security.NewPasswdVerifier(rdb, repo.PasswdErrLimitPrefixKey(), 0, int32(config.Cfg().Service.ErrAttemptLimit)),
		jwtHelper:        security.NewJWT(&config.Cfg().Security.Jwt, rdb),
		repo:             r,
	}
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {

	// 要求
	// 1.校验验证码, 并清空验证码
	// 2.对密码字段从密文转换成明文
	// 3.是否禁用（包括已注销的）
	// 4.校验密码
	// 5.登录次数（防止被暴力破解），登录次数超过5次直接禁用
	// 6.生成新token, 并更新token
	// 7.更新登录日志

	// 1.支持多平台， 不同平台有不同的 token
	// 2.每一个平台只能登录同时在线一个

	// 1. 校验验证码，并清空验证码
	if !s.repo.Captcha.Verify(req.VerificationKey, req.VerificationCode, true) {
		return nil, errs.Status(ctx, errs.Err_CaptchaValidFailed)
	}

	// 2. 对密码字段从密文转换成明文
	passwd, err := s.aesDecrypt(req.Password)
	if err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt",
			"secret", config.Cfg().Security.Ciphertext.CipherKey, // 自动脱敏打印
			"password", req.Password,
			"err", err,
		)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 3. 查询用户信息并检查状态
	user, err := s.repo.User.FindOneUserByUsername(ctx, req.Username)
	if err != nil {
		// 用户名不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.passwdVerifier.Incr(ctx); err != nil {
				slog.ErrorContext(ctx, "s.passwdVerifier.Incr", "err", err)
			}
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

	// 4. 校验密码
	if _, err := s.passwdVerifier.BcryptVerifyWithCount(ctx, string(user.Password), passwd); err != nil {
		// 5. 登录次数（防止被暴力破解），登录次数超过5次直接禁用
		if errors.Is(err, security.ErrPasswdLimit) {
			s.repo.User.UpdateStatus(ctx, user.ID, int32(v1.UserStatus_banned))
		} else {
			slog.ErrorContext(ctx, "s.passwdVerifier.BcryptVerifyWithCount", "err", err)
		}
		return nil, errs.Status(ctx, errs.Err_UsernameOrPasswordErr)
	}

	// 7. 生成新token, 并更新token
	accessToken, refreshToken, err := s.jwtHelper.Generate(ctx, &security.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: user.Username},
		Data: v1.UserTokenInfo{
			UserId:   user.ID,
			Username: user.Username,
			Role:     v1.Role(user.Role),
		},
	})
	if err != nil {
		return nil, errs.Status(ctx, errs.Err_LoginFailed)
	}

	// 7. 更新登录日志
	loginLog := &model.UserLoginSuccessLog{
		UserID:            user.ID,
		Username:          user.Username,
		LoginTime:         time.Now(),
		ClientIP:          "", // TODO 获取客户端IP
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
		TokenType:    "Bearer", // TODO 获取登陆平台
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
	passwd, err := s.aesDecrypt(request.Password)
	if err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt",
			"secret", config.Cfg().Security.Ciphertext.CipherKey, // 自动脱敏打印
			"password", request.Password,
			"err", err,
		)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 3. 密码强度校验
	if level, err := s.passwordStrength.Validate(passwd); err != nil {
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
	hashedPassword := security.Ciphertext(passwd)
	hashedPassword.BcryptHash()

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

	// 要求
	// 1. 对密码字段从密文转换成明文 (新密码和旧密码)
	// 2. 新旧密码不能相同
	// 3. 新密码强度校验
	// 4. 校验用户密码（旧密码）
	// 5. 密码字段加密
	// 6. 更新用户密码
	// 7. 清空用户密码错误次数
	// 8. 清空所有平台的token

	// 1. 对密码字段从密文转换成明文 (新密码和旧密码)
	cipherKey := config.Cfg().Security.Ciphertext.CipherKey
	oloPasswd, err := s.aesDecrypt(request.OldPassword)
	if err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt", "secret", cipherKey, "old password", request.OldPassword, "err", err)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}
	newPasswd, err := s.aesDecrypt(request.NewPassword)
	if err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt", "secret", cipherKey, "new password", request.NewPassword, "err", err)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 2. 新旧密码不能相同
	if oloPasswd == newPasswd {
		return nil, errs.Status(ctx, errs.Err_NewPasswordSameAsOld)
	}

	// 3. 新密码强度校验
	if level, err := s.passwordStrength.Validate(newPasswd); err != nil {
		st, _ := status.New(codes.Code(errs.Err_InvalidPasswordFormat), errs.GetErrorMessage(ctx, errs.Err_InvalidPasswordFormat)).
			WithDetails(&v1.PasswordStrengthError{Level: int32(level)})
		return nil, st.Err()
	}

	// 4. 校验用户密码（旧密码）
	user, err := s.repo.User.FindOneUserByUsername(ctx, request.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Status(ctx, errs.Err_UserNotFound)
		}
		slog.ErrorContext(ctx, "s.repo.UserRepo.FindOneUserByUsername", "username", request.Username, "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}
	if !s.passwdVerifier.BcryptVerify(string(user.Password), oloPasswd) {
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 5. 密码字段加密
	hashedPassword := security.Ciphertext(newPasswd)
	hashedPassword.BcryptHash()

	// 6. 更新用户密码
	if err := s.repo.User.UpdatePassword(ctx, user.ID, string(hashedPassword)); err != nil {
		slog.ErrorContext(ctx, "s.repo.UserRepo.UpdatePassword", "err", err, "user", user)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}

	// 7. TODO 清空用户密码错误次数
	// if err := s.passwdVerifier.Clear(ctx, user.Username); err != nil {
	// 	slog.ErrorContext(ctx, "s.passwdVerifier.Clear", "err", err, "user", user)
	// 	return nil, errs.Status(ctx, errs.Err_SystemBusy)
	// }

	// 8. 清空所有平台的token
	if err := s.removeAllTokens(ctx, user.Username); err != nil {
		slog.ErrorContext(ctx, "s.removeAllTokens", "err", err, "username", user.Username)
	}

	return &emptypb.Empty{}, nil
}

// 生成验证码
func (s *AuthService) SendVerificationCode(ctx context.Context, _ *emptypb.Empty) (*v1.SendVerificationCodeResponse, error) {

	// 要求
	// 1. 生成验证码Key
	// 2. 生成验证码图片

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

	// 要求
	// 1. 校验并解析refreshToken， 获取 username， platform
	// 2. 获取用户信息
	// 3. 生成新的token，并更新token

	// 1. 校验并解析refreshToken， 获取 username， platform
	claims, jwtErr := s.jwtHelper.Parse(request.RefreshToken)
	if jwtErr != nil {
		var err error
		switch {
		case errors.Is(jwtErr, jwt.ErrTokenMalformed):
			err = errs.Status(ctx, errs.Err_TokenMalformed)
		case errors.Is(jwtErr, jwt.ErrTokenUnverifiable):
			err = errs.Status(ctx, errs.Err_TokenUnverifiable)
		case errors.Is(jwtErr, jwt.ErrTokenExpired):
			err = errs.Status(ctx, errs.Err_TokenExpired)
		case errors.Is(jwtErr, jwt.ErrTokenInvalidSubject):
			err = errs.Status(ctx, errs.Err_TokenInvalidSubject)
		default:
			err = errs.Status(ctx, errs.Err_TokenUnverifiable)
		}
		slog.ErrorContext(ctx, "s.jwtHelper.Parse", "err", jwtErr)
		return nil, err
	}
	// 2. 获取用户信息
	user, err := s.repo.User.FindOneUserByUsername(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Status(ctx, errs.Err_UserNotFound)
		}
		slog.ErrorContext(ctx, "s.repo.UserRepo.FindOneUserByUsername", "username", claims.Subject, "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}
	// 3. 生成新的token，并更新token
	accessToken, refreshToken, err := s.jwtHelper.Generate(ctx, &security.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: user.Username},
		Data: v1.UserTokenInfo{
			UserId:   user.ID,
			Username: user.Username,
			Role:     v1.Role(user.Role),
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, "s.jwtHelper.Generate", "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}
	return &v1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer", // TODO 获取登陆平台
		ExpiresIn:    config.Cfg().Basic.Security.Jwt.AccessTokenExpired.TimeDuration().Microseconds(),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	// 要求
	// 1. 获取当前用户信息（username， platform）
	// 2. 移除当前平台的token

	if err := s.jwtHelper.RemoveToken(ctx, ""); err != nil {
		slog.ErrorContext(ctx, "s.jwtHelper.RemoveToken", "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthService) Deregister(ctx context.Context, request *v1.DeregisterRequest) (*emptypb.Empty, error) {

	// 要求
	// 1. 对密码字段从密文转换成明文
	// 2. 校验用户密码
	// 3. 更改用户状态为注销, 并添加注销请求记录
	// 4. 注销用户的所有token

	// 1. 对密码字段从密文转换成明文
	passwd, err := s.aesDecrypt(request.Password)
	if err != nil {
		slog.ErrorContext(ctx, "ciphertext.Decrypt",
			"secret", config.Cfg().Security.Ciphertext.CipherKey, // 自动脱敏打印
			"password", request.Password,
			"err", err,
		)
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 2.验证用户密码
	user, err := s.repo.User.FindOneUserById(ctx, request.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Status(ctx, errs.Err_UserNotFound)
		}
		slog.ErrorContext(ctx, "s.repo.UserRepo.FindOneUserById", "userId", request.UserId, "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}
	if !s.passwdVerifier.BcryptVerify(string(user.Password), passwd) {
		return nil, errs.Status(ctx, errs.Err_InvalidPassword)
	}

	// 3. 更改用户状态为注销, 并添加注销请求记录
	if err := s.repo.User.DeletedUser(ctx, request.UserId); err != nil {
		slog.ErrorContext(ctx, "s.repo.UserRepo.DeleteUser", "userId", request.UserId, "err", err)
		return nil, errs.Status(ctx, errs.Err_SystemBusy)
	}

	// 4. 注销用户的所有token
	if err = s.removeAllTokens(ctx, user.Username); err != nil {
		slog.ErrorContext(ctx, "s.jwtHelper.Deregister", "username", user.Username, "err", err)
		return nil, nil
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthService) removeAllTokens(ctx context.Context, username string) error {
	// TODO 获取多平台
	return s.jwtHelper.RemoveToken(ctx, username)
}

func (s *AuthService) aesDecrypt(ciphertext string) (string, error) {
	if config.Cfg().Security.Ciphertext.IsCiphertext { // 前端是否加密
		return ucrypto.AESDecrypt(ciphertext, string(config.Cfg().Security.Ciphertext.CipherKey))
	}
	return ciphertext, nil
}
