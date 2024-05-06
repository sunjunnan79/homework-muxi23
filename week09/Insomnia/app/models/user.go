package models

import (
	. "Insomnia/app/utility/tool"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid     string `gorm:"size:64;not null;unique"`
	Email    string `gorm:"size:255;not null;unique"`
	Password string `gorm:"size:64;not null"`
	Avatar   string `gorm:"not null"`
}

//--------通用部分------------

// ExistUP 检查email和password是否正确
func ExistUP(email string, password string) bool {
	var user User
	Db.Table("users").First(&user, "email = ? ", email)
	return password == user.Password
}

// ExistEmail 检查email是否存在
func ExistEmail(email string) bool {
	return Db.Table("users").Where("email = ?", email).Error == nil
}

// FindByEmail 方法用于根据邮箱查询用户记录
func FindByEmail(email string) (user User, err error) {
	err = Db.Table("users").First(&user, "email = ? ", email).Error
	return
}

// FindByUuid 方法用于根据Uuid查询用户记录
func FindByUuid(uuid string) (user User, err error) {
	err = Db.Table("users").Where("uuid = ? ", uuid).First(&user).Error
	return
}

// Update 方法用于更新数据库中用户的信息
func (user *User) Update() (err error) {
	//更新user记录
	return Db.Model(user).Updates(map[string]interface{}{"password": user.Password, "email": user.Email, "avatar": user.Avatar}).Error
}

//--------auth部分------------

// Create 方法用于创建新用户并将用户信息保存到数据库中
func (user *User) Create() (err error) {
	//生成用户的UUID
	uuid := CreateUuid()
	user.Uuid = uuid
	//创建新user
	return Db.Create(user).Error
}

//--------task部分------------

// Delete 方法用于从数据库中删除用户
func (user *User) Delete() (err error) {
	//删除user记录
	result := Db.Delete(user)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

// UserDeleteAll 方法用于删除数据库中的"所有"用户记录!
func UserDeleteAll() (err error) {
	//删除所有user记录
	result := Db.Delete(&User{})
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// Users 方法用于获取所有用户的信息
func Users() (users []User, err error) {
	//查询所有user用户记录
	result := Db.Find(&users)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

// UserByName 方法用于根据邮箱查询用户记录
func UserByName(name string) (user User, err error) {
	//根据email查询用户记录
	result := Db.Where("name = ?", name).First(&user)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

// CreatePost 创建一个新的帖子回复
//func (user *User) CreatePost(thread models.Thread, Body string) (post models.Post, err error) {
//	uuid := CreateUuid()
//	now := time.Now()
//	post = models.Post{
//		Uuid:      uuid,
//		Body:      Body,
//		UserId:    user.ID,
//		ThreadId:  thread.ID,
//		CreatedAt: now,
//	}
//	result := models.Db.Create(&post)
//	if result.Error != nil {
//		err = result.Error
//		return
//	}
//	return
//}
