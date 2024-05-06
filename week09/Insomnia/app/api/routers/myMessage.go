package routers

import (
	"Insomnia/app/infrastructure/middlewares"
)

func (r *router) useMyMessage() {
	useMyMessage := r.Group("myMessage")
	useMyMessage.POST("/getLikeMessage", middlewares.UseJwt(), r.myMessage.GetLikeMessage)
	useMyMessage.POST("/checkLikeMessage", middlewares.UseJwt(), r.myMessage.CheckLikeMessage)
	useMyMessage.POST("/getPostMessage", middlewares.UseJwt(), r.myMessage.GetPostMessage)
	useMyMessage.POST("/checkPostMessage", middlewares.UseJwt(), r.myMessage.CheckPostMessage)
}
