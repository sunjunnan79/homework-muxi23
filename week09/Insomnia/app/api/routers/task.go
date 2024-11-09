package routers

import (
	"Insomnia/app/infrastructure/middlewares"
)

// useTask 获取当日数据的路由
func (r *router) useTask() {
	taskRouter := r.Group("task")
	taskRouter.POST("/upTask", middlewares.UseJwt(), r.task.UpTask)
	taskRouter.POST("/getAllTask", middlewares.UseJwt(), r.task.GetAllTask)
}
