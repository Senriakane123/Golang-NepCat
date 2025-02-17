package DBModel

// 管理员用户表
type AdminUser struct {
	ID    int `gorm:"primaryKey;autoIncrement;column:ID"`
	QQNum int `gorm:"type:int;not null;column:QQNum"`
}
