package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tiktok/config"
)

// MysqlHandle global database handle
var MysqlHandle *gorm.DB

func init() {
	db, err := gorm.Open(mysql.Open(config.Mysql.DSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败！"))
	}
	// 自动迁移
	db.AutoMigrate(&UserModel{}, &VideoModel{}, &Favorite{})
	MysqlHandle = db.Debug()
}
