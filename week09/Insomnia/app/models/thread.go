package models

import (
	"Insomnia/app/api/request"
	. "Insomnia/app/utility/tool"
	"encoding/json"
	"gorm.io/gorm"
	"log"
)

type Thread struct {
	gorm.Model
	TUuid  string `gorm:"size:64;not null;unique"`
	Topic  string `gorm:"not null"`
	Title  string `gorm:"size:64;not null;"`
	Uuid   string `gorm:"not null"`
	Likes  uint
	Body   string `gorm:"not null"`
	Number uint
	Images string `gorm:"type:json"`
}

// CreateThread 方法创建一个新的帖子
func CreateThread(UuiD string, ct request.CreateThreadReq) (thread Thread, err error) {
	//生成会话的TUuid
	tUuid := CreateUuid()
	jsonData, err := json.Marshal(ct.Images)
	if err != nil {
		log.Fatal("序列化数据失败:", err)
	}
	thread = Thread{
		Model:  gorm.Model{},
		TUuid:  tUuid,
		Topic:  ct.Topic,
		Uuid:   UuiD,
		Body:   ct.Body,
		Likes:  0,
		Number: 0,
		Title:  ct.Title,
		Images: string(jsonData),
	}
	result := Db.Create(&thread)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// DestroyThread 删除指定的帖子
func DestroyThread(tUuid string) (err error) {
	// 开始事务
	tx := Db.Begin()
	err = Db.Table("threads").Where("t_uuid = ? ", tUuid).Delete(&Thread{}).Error
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	err = Db.Table("posts").Where("t_uuid = ? ", tUuid).Delete(&Post{}).Error
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	err = Db.Table("re_posts").Where("t_uuid = ? ", tUuid).Delete(&RePost{}).Error
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	err = DestroyMessageByTUuId(tUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// NumReplies 方法用于获取帖子的回复数量
func (thread *Thread) NumReplies() (count int64) {
	result := Db.Model(&Post{}).Where("t_uuid = ?", thread.ID).Count(&count)
	if result.Error != nil {
		return
	}
	return
}

// Posts 方法用于获取该帖子的所有回复
func (thread *Thread) Posts() (posts []Post, err error) {
	result := Db.Where("t_uuid = ?", thread.TUuid).Find(&posts)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// Threads 方法用于获取数据库中的所有帖子记录
func Threads(topic string) (threads []Thread, err error) {
	result := Db.Where("topic = ?", topic).Order("created_at desc").Find(&threads)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// ThreadByTUUID 用于根据帖子的TUuid查询帖子记录
func ThreadByTUUID(tUuid string) (thread Thread, err error) {
	result := Db.Where("t_uuid = ?", tUuid).First(&thread)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// AuthorByTUUID 获取当前帖子的作者是谁
func AuthorByTUUID(tUuid string) (author string, err error) {
	thread := Thread{}
	result := Db.Where("t_uuid = ?", tUuid).First(&thread)
	if result.Error != nil {
		err = result.Error
		return
	}
	author = thread.Uuid
	return
}

// ThreadByUuId 用于根据用户的UuId 查询帖子记录
func ThreadByUuId(UuId string) (threads []Thread, err error) {
	result := Db.Where("uuid = ?", UuId).First(&threads)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

func UpThreadLikesData(tUuid string, exist bool) error {
	var thread Thread
	result := Db.Where("t_uuid = ?", tUuid).First(&thread)
	if result.Error != nil {
		err := result.Error
		return err
	}
	if exist {
		thread.Likes++
		return Db.Save(&thread).Error
	}
	thread.Likes--
	return Db.Save(&thread).Error
}

func UpThreadNumbers(tUuid string) error {
	var thread Thread
	result := Db.Where("t_uuid = ?", tUuid).First(&thread)
	if result.Error != nil {
		err := result.Error
		return err
	}
	thread.Number++
	return Db.Save(&thread).Error
}
func DownThreadNumbers(tUuid string) error {
	var thread Thread
	result := Db.Where("t_uuid = ?", tUuid).First(&thread)
	if result.Error != nil {
		err := result.Error
		return err
	}
	thread.Number--
	return Db.Save(&thread).Error
}
