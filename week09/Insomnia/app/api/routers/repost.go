package routers

import (
	"Insomnia/app/infrastructure/middlewares"
)

// useRePost re回复的路由
func (r *router) useRePost() {
	rePostRouter := r.Group("rePost")
	rePostRouter.POST("/createRePost", middlewares.UseJwt(), r.repost.CreateRePost)
	rePostRouter.POST("/destroyRePost", middlewares.UseJwt(), r.repost.DestroyRePost)
	rePostRouter.POST("/readRePost", middlewares.UseJwt(), r.repost.ReadRePost)
	rePostRouter.POST("/getRePosts", middlewares.UseJwt(), r.repost.GetRePosts)
	rePostRouter.POST("/likeRePost", middlewares.UseJwt(), r.repost.LikeRePost)
}
