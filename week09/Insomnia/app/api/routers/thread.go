package routers

import (
	"Insomnia/app/infrastructure/middlewares"
)

// useThread 帖子的路由
func (r *router) useThread() {
	threadRouter := r.Group("thread")
	threadRouter.POST("/createThread", middlewares.UseJwt(), r.thread.CreateThread)
	threadRouter.POST("/getMyThreads", middlewares.UseJwt(), r.thread.GetMyThreads)
	threadRouter.POST("/destroyThread", middlewares.UseJwt(), r.thread.DestroyThread)
	threadRouter.POST("/readThread", middlewares.UseJwt(), r.thread.ReadThread)
	threadRouter.POST("/getThreads", middlewares.UseJwt(), r.thread.GetThreads)
	threadRouter.POST("/likeThread", middlewares.UseJwt(), r.thread.LikeThread)
}
