package routers

import (
	"Insomnia/app/infrastructure/middlewares"
)

// 通用的路由
func (r *router) useCommon() {
	commonRouter := r.Group("common")
	commonRouter.POST("/sendEmail", r.email.SendEmail)
	commonRouter.POST("/getQNToken", middlewares.UseJwt(), r.tube.GetQNToken)
}
