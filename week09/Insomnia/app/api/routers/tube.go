package routers

import "Insomnia/app/infrastructure/middlewares"

func (r *router) useTube() {
	tubeRouter := r.Group("tube")
	tubeRouter.POST("/getQNToken", middlewares.UseJwt(), r.tube.GetQNToken)
}
