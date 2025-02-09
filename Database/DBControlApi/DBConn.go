package DBControlApi

import "log"

// 数据库连接变量
var Db DBcontrol

// 初始化数据库
func InitDatabase() {
	dsn := "username:password@tcp(127.0.0.1:3306)/example_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
}
