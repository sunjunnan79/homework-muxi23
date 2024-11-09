package routers

import (
	"Insomnia/app/infrastructure/middlewares"
)

func (r *router) useAuth() {
	authRouter := r.Group("auth")
	authRouter.POST("/login", r.auth.Login)
	authRouter.POST("/signup", r.auth.Signup)
	authRouter.POST("/changePassword", r.auth.ChangePassword)
	authRouter.POST("/changeAvatar", middlewares.UseJwt(), r.auth.ChangeAvatar)
	authRouter.POST("/getMyData", middlewares.UseJwt(), r.auth.GetMyData)
}
