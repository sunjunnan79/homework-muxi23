package controller

import (
	. "Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	. "Insomnia/app/infrastructure/helper"
	"Insomnia/app/service"
	"Insomnia/app/utility/tool"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Auth struct{}

var authService *service.AuthService

// Login 用户登录
// @Summary 用户登录接口
// @Description 用户登录接口
// @Tags Auth
// @Accept json
// @Produce json
// @Param email query string true "邮箱"
// @Param password query string true "密码"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/auth/login [post]
func (a *Auth) Login(c *gin.Context) {
	//定义一个Login请求类型的结构体
	req := &LoginReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		Danger(err, "解析请求结构体失败")
		FailMsgData(c, fmt.Sprintf("params invalid error: %v", err), LoginResponse{})
		return
	}

	//调用服务层来获取一个token
	token, err := authService.Login(req.Email, tool.Encrypt(req.Password))
	if err != nil {
		Danger(err, "获取token失败")
		FailMsgData(c, fmt.Sprintf("获取token失败: %v", err), LoginResponse{})
		return
	}

	//返回消息捏
	OkMsgData(c, "登录成功", LoginResponse{Token: token})
}

// Signup 用户注册
// @Summary 用户注册接口
// @Description 用户注册接口
// @Tags Auth
// @Accept json
// @Produce json
// @Param email query string true "邮箱"
// @Param password query string true "密码"
// @Param verificationCode query string true "验证码"
// @Param avatar query string false "头像"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/auth/signup [post]
func (a *Auth) Signup(c *gin.Context) {
	//定义一个Login请求类型的结构体
	sur := &SignupReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(sur); err != nil {
		Danger(err, "无法解析的表单")
		FailMsg(c, fmt.Sprintf("无法解析: %v", err))
		return
	}

	//调用服务层来注册一个账户
	err := authService.Signup(*sur)
	if err != nil {
		Danger(err, "注册时服务器发生错误")
		FailMsg(c, fmt.Sprintf("注册时服务器发生错误: %v", err))
		return
	}

	//返回消息捏
	OkMsg(c, "注册成功!")
}

// ChangePassword 更改密码
// @Summary 更改密码接口
// @Description 更改密码接口
// @Tags Auth
// @Accept json
// @Produce json
// @Param email query string true "邮箱"
// @Param verificationCode query string true "验证码"
// @Param Authorization header string true "jwt验证"
// @Param newPassword query string true "新密码"
// @Success 200 {object} string "密码更改成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/auth/changePassword [post]
func (a *Auth) ChangePassword(c *gin.Context) {
	//定义一个Login请求类型的结构体
	cp := &ChangePasswordReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(cp); err != nil {
		Danger(err, "无法解析的表单")
		FailMsg(c, fmt.Sprintf("无法解析: %v", err))
		return
	}

	//调用服务层来更新密码
	err := authService.ChangePassword(*cp)
	if err != nil {
		Danger(err, "更新密码失败")
		FailMsg(c, fmt.Sprintf("更新密码失败: %v", err))
		return
	}

	//返回消息捏
	OkMsg(c, "更改密码成功!")
}

// ChangeAvatar 更改头像
// @Summary 更改头像接口
// @Description 更改头像接口
// @Tags Auth
// @Accept json
// @Produce json
// @Param newAvatar query string true "新头像"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} string "头像更改成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/v1/auth/changeAvatar [post]
func (a *Auth) ChangeAvatar(c *gin.Context) {
	//定义一个Login请求类型的结构体
	cs := &ChangeAvatarReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(cs); err != nil {
		Danger(err, "无法解析的表单")
		FailMsg(c, fmt.Sprintf("无法解析: %v", err))
		return
	}

	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	//调用服务层来更新头像
	err := authService.ChangeAvatar(*cs, uuid)
	if err != nil {
		Danger(err, "更新头像失败")
		FailMsg(c, fmt.Sprintf("更新头像失败: %v", err))
		return
	}

	//返回消息捏
	OkMsg(c, "更改头像成功!")
}

// GetMyData 获取数据
// @Summary 获取数据
// @Description 获取数据
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetMyDataResponse "获取数据成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/v1/auth/getMyData [post]
func (a *Auth) GetMyData(c *gin.Context) {
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	myData, err := authService.GetMyData(uuid)
	if err != nil {
		Danger(err, "获取用户信息失败")
		FailMsgData(c, fmt.Sprintf("获取用户信息失败: %v", err), myData)
		return
	}
	OkMsgData(c, "获取用户信息成功", myData)
	return
}
