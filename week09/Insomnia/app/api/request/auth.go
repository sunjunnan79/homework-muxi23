package request

// LoginReq 登陆的请求
type LoginReq struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

// SignupReq 注册的请求
type SignupReq struct {
	Email            string `json:"email" form:"email" binding:"required"`
	Password         string `json:"password" form:"password" binding:"required,min=6"`
	VerificationCode string `json:"verificationCode" form:"verificationCode" binding:"required,alphanum"`
	Avatar           string `json:"avatar" form:"avatar" binding:""`
}

// ChangePasswordReq 更换密码的请求
type ChangePasswordReq struct {
	Email            string `json:"email" form:"email" binding:"required"`
	VerificationCode string `json:"verificationCode" form:"verificationCode" binding:"required,alphanum"`
	NewPassword      string `json:"newPassword" form:"newPassword" binding:"required,min=6"`
}

// ChangeAvatarReq 更换头像的请求
type ChangeAvatarReq struct {
	NewAvatar string `json:"newAvatar" form:"newAvatar" binding:"required"`
}
