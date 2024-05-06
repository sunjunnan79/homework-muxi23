package routers

import (
	"Insomnia/app/controller"
	"github.com/gin-gonic/gin"
)

type router struct {
	*gin.RouterGroup
	auth      *controller.Auth
	task      *controller.Task
	thread    *controller.Thread
	email     *controller.SendEmail
	tube      *controller.Tube
	post      *controller.Post
	repost    *controller.RePost
	myMessage *controller.MyMessage
}

func Load(e *gin.Engine) {
	r := &router{
		RouterGroup: &e.RouterGroup,
		auth:        &controller.Auth{},
		task:        &controller.Task{},
		thread:      &controller.Thread{},
		tube:        &controller.Tube{},
		post:        &controller.Post{},
		repost:      &controller.RePost{},
		myMessage:   &controller.MyMessage{},
	}
	r.RouterGroup = r.Group("/api/v1")
	//启用认证的路由
	r.useAuth()
	//启用任务的路由
	r.useTask()
	//启用帖子的路由
	r.useThread()
	//启动回复的路由
	r.usePost()
	//启动re回复的路由
	r.useRePost()
	//启动通用的路由
	r.useCommon()
	//启动消息的路由
	r.useMyMessage()
	//启用Tube
	r.useTube()
}
