package repo

type All struct {
	Captcha             *CaptchaRepo
	User                *UserRepo
	UserDeletionRequest *UserDeletionRequestRepo
	UserLoginSuccessLog *UserLoginSuccessLogRepo
}

func New(
	captcha *CaptchaRepo,
	user *UserRepo,
	userDeletionRequest *UserDeletionRequestRepo,
	userLoginSuccessLog *UserLoginSuccessLogRepo,
) *All {
	return &All{
		Captcha:             captcha,
		User:                user,
		UserDeletionRequest: userDeletionRequest,
		UserLoginSuccessLog: userLoginSuccessLog,
	}
}
