package model

import (
	"fmt"
	"tiktok/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MysqlHandle global database handle
var MysqlHandle *gorm.DB

func init() {
	db, err := gorm.Open(mysql.Open(config.Mysql.DSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败！"))
	}
	//自动迁移,会创建UserModel,VideoModel.
	db.AutoMigrate(&UserModel{}, &VideoModel{}, &FriendListModel{},&MessageModel{})
	MysqlHandle = db.Debug()
}
