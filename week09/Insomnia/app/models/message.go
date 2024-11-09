package models

import (
	"gorm.io/gorm"
)

type PostMessage struct {
	gorm.Model
	TUuid string `gorm:"size:64"`
	Uuid  string `gorm:"size:64"`
	PUuid string `gorm:"size:64"`
	RUuid string `gorm:"size:64"`
	Title string `gorm:"size:64"`
	Check bool   `gorm:"not null"`
	Body  string `gorm:"not null"`
	Likes uint
}

type LikeMessage struct {
	gorm.Model
	TUuid string `gorm:"size:64"`
	Uuid  string `gorm:"size:64"`
	PUuid string `gorm:"size:64"`
	RUuid string `gorm:"size:64"`
	Title string `gorm:"size:64"`
	Check bool   `gorm:"not null"`
	Body  string `gorm:"not null"`
	Likes uint
}

func DeleteLikeMessage(lm LikeMessage) error {
	return Db.Table("like_messages").Unscoped().Where("t_uuid = ? AND p_uuid = ? AND r_uuid = ? AND uuid = ?", lm.TUuid, lm.RUuid, lm.PUuid, lm.Uuid).Delete(&lm).Error
}

// CreatePostMessage 方法创建一个新的消息提醒
func CreatePostMessage(pm PostMessage) (err error) {
	result := Db.Create(&pm)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// CheckPostMessage 删除指定的消息提醒
func CheckPostMessage(id string) (err error) {
	var pm PostMessage
	err = Db.Table("post_messages").Where("id = ? ", id).First(&pm).Error
	if err != nil {
		return
	}

	if pm.Check {
		pm.Check = false
		return Db.Save(&pm).Error
	}
	return
}

// PostMessageByUuId 用于根据用户的UuId查询消息提醒
func PostMessageByUuId(UuId string) (pm []PostMessage, err error) {
	result := Db.Where("uuid = ?", UuId).Order("created_at desc").Find(&pm)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// CreateLikeMessage 方法创建一个新的消息提醒
func CreateLikeMessage(lm LikeMessage) (err error) {
	result := Db.Create(&lm)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// CheckLikeMessage 检查指定的消息提醒
func CheckLikeMessage(id string) (err error) {
	var lm LikeMessage
	err = Db.Table("like_messages").Where("id = ? ", id).First(&lm).Error
	if err != nil {
		return
	}

	if !lm.Check {
		lm.Check = true
		return Db.Save(&lm).Error
	}
	return
}

// LikeMessageByUuId 用于根据用户的UuId查询消息提醒
func LikeMessageByUuId(UuId string) (lm []LikeMessage, err error) {
	err = Db.Table("like_messages").Where("uuid = ?", UuId).Order("created_at desc").Find(&lm).Error
	return
}

// DestroyMessageByTUuId 删除某个帖子下的所有消息提醒
func DestroyMessageByTUuId(tUuId string) (err error) {
	// 开始事务
	tx := Db.Begin()
	err = Db.Table("like_messages").Where("t_uuid = ?", tUuId).Delete(&LikeMessage{}).Error
	if err != nil {
		//事务回滚
		tx.Rollback()
		return
	}
	err = Db.Table("post_messages").Where("t_uuid = ?", tUuId).Delete(&PostMessage{}).Error
	if err != nil {
		//事务回滚
		tx.Rollback()
		return
	}
	//提交事务
	tx.Commit()
	return
}

// DestroyMessageByPUuId 删除某个回复下的所有消息提醒
func DestroyMessageByPUuId(pUuId string) (err error) {
	// 开始事务
	tx := Db.Begin()
	err = Db.Table("like_messages").Where("p_uuid = ?", pUuId).Delete(&LikeMessage{}).Error
	if err != nil {
		//事务回滚
		tx.Rollback()
		return
	}
	err = Db.Table("post_messages").Where("p_uuid = ?", pUuId).Delete(&PostMessage{}).Error
	if err != nil {
		//事务回滚
		tx.Rollback()
		return
	}
	//提交事务
	tx.Commit()
	return
}

// DestroyMessageByRUuId 删除某个re回复下的所有消息提醒
func DestroyMessageByRUuId(rUuId string) (err error) {
	// 开始事务
	tx := Db.Begin()
	err = Db.Table("like_messages").Where("r_uuid = ?", rUuId).Delete(&LikeMessage{}).Error
	if err != nil {
		//事务回滚
		tx.Rollback()
		return
	}
	err = Db.Table("post_messages").Where("r_uuid = ?", rUuId).Delete(&PostMessage{}).Error
	if err != nil {
		//事务回滚
		tx.Rollback()
		return
	}
	//提交事务
	tx.Commit()
	return
}
