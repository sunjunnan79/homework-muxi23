package controller

import (
	"Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	. "Insomnia/app/infrastructure/helper"
	"Insomnia/app/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Task struct{}

var taskService *service.TaskService

// UpTask 更新今日数据(星期),将今日完成的个数加一
// @Summary 更新数据
// @Description 更新数据
// @Tags Task
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetMyDataResponse "更新数据成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/v1/task/upTask [post]
func (t *Task) UpTask(c *gin.Context) {
	//获取当前用户的Uuid
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	//定义一个Login请求类型的结构体
	req := &request.UpdateTaskNumber{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		Danger(err, "解析请求结构体失败")
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	err := taskService.UpTask(uuid, req.Num)
	if err != nil {
		FailMsg(c, err.Error())
		return
	}

	OkMsg(c, "更新数据成功!")
	return
}

// GetAllTask 获取本周数据(星期)
// @Summary 获取数据
// @Description 获取数据
// @Tags Task
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} []GetTaskResponse "获取数据成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/v1/task/getAllTask [post]
func (t *Task) GetAllTask(c *gin.Context) {
	//获取当前用户的Uuid
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	err, allTask := taskService.GetAllTask(uuid)
	if err != nil {
		Danger(err, "获取本周数据失败")
		FailMsgData(c, err.Error(), []int{})
		return
	}
	OkMsgData(c, "获取本周数据成功", allTask)
	return
}
