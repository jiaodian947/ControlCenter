package auth

type LoginForm struct {
	UserName string `valid:"Required"`
	Password string `form:"type(password)" valid:"Required"`
	Captcha  string `valid:"Required"`
	Remember bool
}

func (form *LoginForm) Labels() map[string]string {
	return map[string]string{
		"UserName": "auth.username_or_email",
		"Password": "auth.login_password",
		"Remember": "auth.login_remember_me",
	}
}

type CreateForm struct {
	UserName string `valid:"Required"`
	Email    string `valid:"Required"`
	Captcha  string `valid:"Required"`
	PassWord string `form:"type(password)" valid:"Required"`
}
