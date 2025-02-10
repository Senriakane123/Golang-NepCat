package GameDatamodel

// 宠物表
type Pet struct {
	ID                  int     `gorm:"primaryKey;autoIncrement"`
	Name                string  `gorm:"type:varchar(100);not null"`
	Type                string  `gorm:"type:varchar(50);not null"`
	Skill               string  `gorm:"type:varchar(255)"`
	HealthGrowthFactor  float32 `gorm:"type:float"`
	AtkGrowthFactor     float32 `gorm:"type:float"`
	DefenseGrowthFactor float32 `gorm:"type:float"`
	EnergyGrowthFactor  float32 `gorm:"type:float"`
	BaseHealth          int     `gorm:"type:int;not null"`
	BaseAtk             int     `gorm:"type:int;not null"`
	BaseDef             int     `gorm:"type:int;not null"`
	BaseEnergy          int     `gorm:"type:int;not null"`
}

// 用户信息表
type UserInfo struct {
	ID      int               `gorm:"primaryKey;autoIncrement"`
	QQNum   int               `gorm:"type:int;not null"`
	Name    string            `gorm:"type:varchar(100);not null"`
	Item    string            `gorm:"type:varchar(255)"`
	PetInfo []PersonalPetInfo `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// 用户宠物信息表
type PersonalPetInfo struct {
	ID       int                `gorm:"primaryKey;autoIncrement"`
	UserID   int                `gorm:"index;not null"` // 外键，关联 UserInfo
	PetId    int                `gorm:"not null"`       // 外键，关联 Pet
	Petlevel int                `gorm:"type:int;not null"`
	Exp      int                `gorm:"type:int;not null"`
	Skill    []PesonalSkillList `gorm:"foreignKey:PersonalPetID;constraint:OnDelete:CASCADE"`
}

// 技能表
type PesonalSkillList struct {
	ID            int `gorm:"primaryKey;autoIncrement"`
	SkillID       int `gorm:"not null"` //关联AllSkillList里的ID
	PersonalPetID int `gorm:"index"`    // 外键，关联 PersonalPetInfo
}

type AllSkillList struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	SkillName string `gorm:"type:varchar(100);not null"`
	Des       string `gorm:"type:varchar(255)"`
	PetID     []int  `gorm:"type:int;not null"` //只有这个数组里对于petid里的宠物可以学习这个技能
}
