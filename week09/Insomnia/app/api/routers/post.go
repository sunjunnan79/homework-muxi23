package routers

import (
	"Insomnia/app/infrastructure/middlewares"
)

// usePost 回复的路由
func (r *router) usePost() {
	postRouter := r.Group("post")
	postRouter.POST("/createPost", middlewares.UseJwt(), r.post.CreatePost)
	postRouter.POST("/destroyPost", middlewares.UseJwt(), r.post.DestroyPost)
	postRouter.POST("/readPost", middlewares.UseJwt(), r.post.ReadPost)
	postRouter.POST("/getPosts", middlewares.UseJwt(), r.post.GetPosts)
	postRouter.POST("/likePost", middlewares.UseJwt(), r.post.LikePost)
}
