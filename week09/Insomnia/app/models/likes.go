package models

import (
	"gorm.io/gorm"
)

type Likes struct {
	gorm.Model
	Uid   string `gorm:"size:64;not null"`
	Uuid  string `gorm:"size:64;not null"`
	Exist bool   `gorm:"not null"`
}

func ChangeLike(uid, uuid string) (bool, error) {
	// 在数据库中查找是否存在对应的记录
	like := Likes{}
	result := Db.Where("uid = ? AND uuid = ?", uid, uuid).First(&like)
	if result.Error != nil {
		return false, result.Error
	}

	// 更改点赞状态
	like.Exist = !like.Exist

	// 保存更改
	result = Db.Save(&like)
	if result.Error != nil {
		return false, result.Error
	}

	return like.Exist, nil
}

func CheckLike(uid string, uuid string) (bool, error) {
	// 在数据库中查找是否存在对应的记录
	like := Likes{
		Uid:   uid,
		Uuid:  uuid,
		Exist: false,
	}

	result := Db.Where("uid = ? AND uuid = ?", uid, uuid).FirstOrCreate(&like)
	if result.Error != nil {
		return false, result.Error
	}
	return false, nil
}
