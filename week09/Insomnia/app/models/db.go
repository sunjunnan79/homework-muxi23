package models

import (
	. "Insomnia/app/infrastructure/config"
	"Insomnia/app/utility/tube"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Db 连接数据库的指针
var Db *gorm.DB

// 数据库连接启动!
func init() {
	var err error
	config := LoadConfig()

	//打开数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", config.Db.User, config.Db.Password, config.Db.Server, config.Db.Port, config.Db.Database, config.Db.Config)

	//检测数据库连接是否正常
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 创建用户数据库表
	err = Db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	// 创建帖子数据库表
	err = Db.AutoMigrate(&Thread{})
	if err != nil {
		log.Fatal(err)
	}

	// 创建任务数据库表
	err = Db.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal(err)
	}

	// 创建回复数据库表
	err = Db.AutoMigrate(&Post{})
	if err != nil {
		log.Fatal(err)
	}

	// 创建re回复数据库表
	err = Db.AutoMigrate(&RePost{})
	if err != nil {
		log.Fatal(err)
	}
	//创建点赞消息数据库
	err = Db.AutoMigrate(&LikeMessage{})
	if err != nil {
		log.Fatal(err)
	}
	//创建评论消息数据库
	err = Db.AutoMigrate(&PostMessage{})
	if err != nil {
		log.Fatal(err)
	}
	// 创建点赞数据库表
	err = Db.AutoMigrate(&Likes{})
	if err != nil {
		log.Fatal(err)
	}

	//加载七牛云的参数
	tube.LoadQiniu()
	return
}
