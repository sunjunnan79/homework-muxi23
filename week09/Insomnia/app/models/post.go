package models

import (
	"Insomnia/app/api/request"
	. "Insomnia/app/utility/tool"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	TUuid  string `gorm:"size:64;not null"`
	Uuid   string `gorm:"size:64;not null"`
	PUuid  string `gorm:"size:64;not null;unique"`
	Body   string `gorm:"not null"`
	Number uint
	Likes  uint
}

// CreatePost 方法创建一个新的回复
func CreatePost(UuID string, cp request.CreatePostReq) (post Post, err error) {
	//生成会话的TUuid
	pUuid := CreateUuid()
	post = Post{
		TUuid:  cp.TUuid,
		Uuid:   UuID,
		PUuid:  pUuid,
		Body:   cp.Body,
		Number: 0,
		Likes:  0,
	}
	result := Db.Create(&post)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// DestroyPost 删除指定的回复
func DestroyPost(pUuid string) (err error) {
	// 开始事务
	tx := Db.Begin()
	err = Db.Table("posts").Where("p_uuid = ? ", pUuid).Delete(&Post{}).Error
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	err = Db.Table("re_posts").Where("p_uuid = ? ", pUuid).Delete(&RePost{}).Error
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	err = DestroyMessageByPUuId(pUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// Posts 方法用于获取帖子的所有回复
func Posts(tUuid string) (posts []Post, err error) {
	result := Db.Table("posts").Where("t_uuid = ?", tUuid).Find(&posts)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// PostByPUUID 用于根据回复的PUuid查询帖子记录
func PostByPUUID(pUuid string) (post Post, err error) {
	result := Db.Table("posts").Where("p_uuid = ?", pUuid).First(&post)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// PostByUuId 用于根据用户的UuId 查询回复记录
func PostByUuId(UuId string) (posts []Post, err error) {
	result := Db.Where("uuid = ?", UuId).First(&posts)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// AuthorByPUUID 获取当前帖子的作者是谁
func AuthorByPUUID(pUuid string) (author string, err error) {
	post := Post{}
	result := Db.Where("p_uuid = ?", pUuid).First(&post)
	if result.Error != nil {
		err = result.Error
		return
	}
	author = post.Uuid
	return
}

func UpPostLikesData(pUuid string, exist bool) error {
	var post Post
	result := Db.Where("p_uuid = ?", pUuid).First(&post)
	if result.Error != nil {
		err := result.Error
		return err
	}
	if exist {
		post.Likes++
		return Db.Save(&post).Error
	}
	post.Likes--
	return Db.Save(&post).Error
}
