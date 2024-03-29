package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mock_douyin_project/config"
)

var DB *gorm.DB

func InitDB() {

	var err error
	DB, err = gorm.Open(mysql.Open(config.DBConnectString()), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
		//Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&UserInfo{}, &Video{}, &Comment{}, &UserLogin{}, &Message{}, &UserRelation{})
	if err != nil {
		panic(err)
	}
}
