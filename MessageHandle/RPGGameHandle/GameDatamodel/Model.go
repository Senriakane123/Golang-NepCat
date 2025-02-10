package GameDatamodel

// 宠物表
type Pet struct {
	ID                  int     `gorm:"primaryKey;autoIncrement;column:ID"`     // 表中字段名为 id
	Name                string  `gorm:"type:varchar(100);not null;column:Name"` // 表中字段名为 name
	Type                string  `gorm:"type:varchar(50);not null;column:tTpe"`  // 表中字段名为 type
	Skill               string  `gorm:"type:varchar(255);column:Skill"`         // 表中字段名为 skill
	HealthGrowthFactor  float32 `gorm:"type:float;column:HealthGrowthFactor"`   // 表中字段名为 health_growth_factor
	AtkGrowthFactor     float32 `gorm:"type:float;column:AtkGrowthFactor"`      // 表中字段名为 atk_growth_factor
	DefenseGrowthFactor float32 `gorm:"type:float;column:DefenseGrowthFactor"`  // 表中字段名为 defense_growth_factor
	EnergyGrowthFactor  float32 `gorm:"type:float;column:EnergyGrowthFactor"`   // 表中字段名为 energy_growth_factor
	BaseHealth          int     `gorm:"type:int;not null;column:BaseHealth"`    // 表中字段名为 base_health
	BaseAtk             int     `gorm:"type:int;not null;column:BaseAtk"`       // 表中字段名为 base_atk
	BaseDef             int     `gorm:"type:int;not null;column:BaseDef"`       // 表中字段名为 base_def
	BaseEnergy          int     `gorm:"type:int;not null;column:BaseEnergy"`    // 表中字段名为 base_energy
}

// 用户信息表
type UserInfo struct {
	ID      int               `gorm:"primaryKey;autoIncrement"`                      // 主键自增
	QQNum   int64             `gorm:"column:QQNum;type:int;not null"`                // 数据库列名为 qq_num
	Name    string            `gorm:"column:Name;type:varchar(100);not null"`        // 数据库列名为 user_name
	Item    string            `gorm:"column:Item;type:varchar(255);default:''"`      // 数据库列名为 item_info
	PetInfo []PersonalPetInfo `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // 外键字段，删除时级联
}

// 用户宠物信息表
type PersonalPetInfo struct {
	ID       int                `gorm:"primaryKey;autoIncrement;column:ID"` // 表中字段名为 id
	UserID   int64              `gorm:"index;not null;column:UserID"`       // 表中字段名为 user_id
	PetId    int64              `gorm:"not null;column:PetId"`              // 表中字段名为 pet_id
	Petlevel int                `gorm:"type:int;not null;column:Petlevel"`  // 表中字段名为 pet_level
	Exp      int                `gorm:"type:int;not null;column:Exp"`       // 表中字段名为 exp
	QQNum    int                `gorm:"type:int;not null;column:QQNum"`     // 表中字段名为 exp
	Skill    []PesonalSkillList `gorm:"foreignKey:PersonalPetID;constraint:OnDelete:CASCADE"`
}

// 技能表
type PesonalSkillList struct {
	ID            int `gorm:"primaryKey;autoIncrement;column:ID"` // 表中字段名为 id
	SkillID       int `gorm:"not null;column:SkillID"`            // 表中字段名为 skill_id
	PersonalPetID int `gorm:"index;column:PersonalPetID"`         // 表中字段名为 personal_pet_id
}

type AllSkillList struct {
	ID        int    `gorm:"primaryKey;autoIncrement;column:ID"`          // 表中字段名为 id
	SkillName string `gorm:"type:varchar(100);not null;column:SkillName"` // 表中字段名为 skill_name
	Des       string `gorm:"type:varchar(255);column:Des"`                // 表中字段名为 des
	PetID     []int  `gorm:"type:int;not null;column:PetID"`              // 表中字段名为 pet_id
}
