package service

import (
	. "Insomnia/app/models"
	"fmt"
)

type TaskService struct{}

// UpTask /POST /task1/upTask
// 更新任务数据
func (t *TaskService) UpTask(uuid string, num uint) (err error) {

	//获取当日任务数据
	task, err := GetAllTaskByUuid(uuid)
	if err != nil {
		return fmt.Errorf("用户任务数据查询错误:%v", err)
	}

	//更新任务数据
	err = task.Update(num)
	if err != nil {
		return fmt.Errorf("用户任务数据更新失败%v", err)
	}

	return nil
}

// GetAllTask /GET /task1/getAllTask
// 获取本周的数据
func (t *TaskService) GetAllTask(uuid string) (error, []uint) {
	tasks := make([]uint, 7)

	//获取本周的任务数据
	task, err := GetAllTaskByUuid(uuid)
	if err != nil {
		return fmt.Errorf("用户任务数据查询错误:%v", err), tasks
	}

	tasks[0] = task.Monday
	tasks[1] = task.Tuesday
	tasks[2] = task.Wednesday
	tasks[3] = task.Thursday
	tasks[4] = task.Friday
	tasks[5] = task.Saturday
	tasks[6] = task.Sunday
	return nil, tasks
}
