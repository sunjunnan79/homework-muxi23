package service

import (
	. "Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	"Insomnia/app/infrastructure/Email"
	"Insomnia/app/infrastructure/jwt"
	"Insomnia/app/models"
	"Insomnia/app/utility/tool"
	"fmt"
)

type AuthService struct{}

// Signup 注册新用户
func (a *AuthService) Signup(signup SignupReq) (err error) {
	//获取用户信息
	eck := Email.CheckEmail{
		Email:            signup.Email,
		VerificationCode: signup.VerificationCode,
	}

	//检查验证码是否正确
	err = eck.CheckVerificationCode()
	if err != nil {
		return err
	}

	//检查邮箱是否已经被注册
	if !models.ExistEmail(signup.Email) {
		return fmt.Errorf("邮箱已经被注册了")
	}

	//创建一个user
	user1 := models.User{
		Email:    signup.Email,
		Password: tool.Encrypt(signup.Password),
		Avatar:   signup.Avatar,
	}

	//创建一个新的账户
	err = user1.Create()
	if err != nil {
		return err
	}
	return nil
}

// Login 登陆
func (a *AuthService) Login(email string, password string) (string, error) {
	if !models.ExistUP(email, password) {
		return "", fmt.Errorf("邮箱或密码错误")
	}

	result, err := models.FindByEmail(email)
	if err != nil {
		fmt.Println("你不该看见这个的")
		return "", err
	}

	return jwt.SignToken(result.Uuid)
}

// ChangePassword 更新密码
func (a *AuthService) ChangePassword(cp ChangePasswordReq) (err error) {
	//获取用户信息
	eck := Email.CheckEmail{
		Email:            cp.Email,
		VerificationCode: cp.VerificationCode,
	}

	//检查验证码是否正确
	err = eck.CheckVerificationCode()
	if err != nil {
		return err
	}
	//获取用户信息
	user, err := models.FindByEmail(cp.Email)
	if err != nil {
		return err
	}

	//获取新密码
	user.Password = tool.Encrypt(cp.NewPassword)

	//更新用户信息
	err = user.Update()
	if err != nil {
		return err
	}

	return nil
}

// ChangeAvatar 更换头像
func (a *AuthService) ChangeAvatar(cs ChangeAvatarReq, Uuid string) (err error) {
	//获取用户信息
	user, err := models.FindByUuid(Uuid)
	if err != nil {
		return err
	}

	//获取新头像
	user.Avatar = cs.NewAvatar

	//更新用户信息
	err = user.Update()
	if err != nil {
		return err
	}

	return nil
}

// GetMyData 获取数据
func (a *AuthService) GetMyData(Uuid string) (mm GetMyDataResponse, err error) {
	var amount, likes, getPosts uint
	threads, err := models.ThreadByUuId(Uuid)
	if err != nil {
		return GetMyDataResponse{}, err
	}
	for _, t := range threads {
		likes += t.Likes
		getPosts += t.Number
	}

	posts, err := models.PostByUuId(Uuid)
	if err != nil {
		return GetMyDataResponse{}, err
	}
	for _, t := range threads {
		likes += t.Likes
		getPosts += t.Number
	}

	reposts, err := models.RepostByUuId(Uuid)
	if err != nil {
		return GetMyDataResponse{}, err
	}
	for _, t := range threads {
		likes += t.Likes
	}

	amount = uint(len(threads) + len(posts) + len(reposts))
	mm.MyPost = amount
	mm.Likes = likes
	mm.GetPost = getPosts
	return
}
