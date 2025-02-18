package MusicDBModel

type DBMusicListIDModel struct {
	ID            int    `gorm:"primaryKey;autoIncrement;column:ID"`   // 表中字段名为 id
	MusicListID   int    `gorm:"type:int;not null;column:MusicListID"` // 表中字段名为 pet_level
	MusicListname string `gorm:"column:MusicListname;type:varchar(255);default:''"`
}

type Song struct {
	Name   string
	Singer string
	Album  string
	Time   string
	ID     string
}
