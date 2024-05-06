package controller

import (
	. "Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	. "Insomnia/app/infrastructure/Email"
	. "Insomnia/app/infrastructure/config"
	. "Insomnia/app/infrastructure/helper"
	. "Insomnia/app/utility/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

type SendEmail struct{}

// SendEmail 发送邮件验证码
// @Summary 发送验证码接口
// @Description 发送验证码接口
// @Tags SendEmail
// @Accept json
// @Produce json
// @version 1.0
// @Param email query SendEmailRequest true "邮箱"
// @Success 200 {string} string "发送邮件成功"
// @Failure 404 {string} string "邮箱服务器出错"
// @Failure 500 {string} string "发送邮件失败"
// @Router /api/v1/common/sendEmail [post]
func (e *SendEmail) SendEmail(c *gin.Context) {
	//设置log参数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//获取全局配置
	config := LoadConfig()

	em := email.NewEmail()

	//设置发送方的邮箱,此处可以写自己的邮箱
	em.From = config.Email.UserName + "<" + config.Email.Sender + ">"

	//获取随机验证码
	random := GetRandom()
	//定义一个Login请求类型的结构体
	cp := &SendEmailRequest{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(cp); err != nil {
		Danger(err, "无法解析的表单")
		FailMsg(c, fmt.Sprintf("无法解析: %v", err))
		return
	}
	//创造一个临时的CheckEmail
	Email := CheckEmail{
		Email:            cp.Email,
		VerificationCode: random,
	}

	err := Email.CreateRedis()
	if err != nil {
		FailMsg(c, fmt.Sprintf("服务器出错: %v", err))
		Danger(err, "存储验证码到Redis失败")
		return
	}
	//设置接收方的邮箱
	em.To = []string{Email.Email}

	// 设置主题
	em.Subject = "验证码"

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.Text = []byte(Email.VerificationCode + "(验证码将在5分钟后失效,请不要告诉其他人，并尽快注册。)")

	//设置服务器相关的配置
	err = em.Send(config.Email.Smtp, smtp.PlainAuth("", config.Email.Sender, config.Email.Password, "smtp.qq.com"))
	if err != nil {
		FailMsg(c, fmt.Sprintf("服务器出错: %v", err))
		Danger(err, "邮箱服务器配置失败")
		return
	}

	//提示发送成功
	OkMsg(c, fmt.Sprintf("发送邮件成功!"))
	return
}
