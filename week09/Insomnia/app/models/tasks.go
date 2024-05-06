package models

import (
	"gorm.io/gorm"
	"time"
)

// Task 个人总数据
type Task struct {
	gorm.Model
	Uuid      string
	Monday    uint
	Tuesday   uint
	Wednesday uint
	Thursday  uint
	Friday    uint
	Saturday  uint
	Sunday    uint
}

//Day := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

// Update 方法更新数据
func (task Task) Update(num uint) (err error) {

	//获取当前时间
	currentTime := time.Now()
	switch int(currentTime.Weekday()) {
	case 0:
		task.Sunday = num
	case 1:
		task.Monday = num
	case 2:
		task.Tuesday = num
	case 3:
		task.Wednesday = num
	case 4:
		task.Thursday = num
	case 5:
		task.Friday = num
	case 6:
		task.Saturday = num
	}

	result := Db.Model(task).Table("tasks").Where("uuid =?", task.Uuid).Save(&task)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// GetAllTaskByUuid 方法获取所有的task数据
func GetAllTaskByUuid(uuid string) (task *Task, err error) {
	// 根据uuid查询任务记录
	result := Db.Where("uuid = ?", uuid).First(&task)

	if result.Error != nil {
		task.Sunday = 0
		task.Monday = 0
		task.Tuesday = 0
		task.Wednesday = 0
		task.Thursday = 0
		task.Friday = 0
		task.Saturday = 0
		task.Uuid = uuid

		result = Db.Create(&task)
		if result.Error != nil {
			err = result.Error
			return
		}

		result = Db.Where("uuid = ?", uuid).First(&task)
		if result.Error != nil {
			err = result.Error
			return
		}
	}

	return
}
