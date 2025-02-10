package MessageModel

// 定义基础消息结构
type Message struct {
	PostType      string  `json:"post_type"`
	SelfID        int64   `json:"self_id"`
	UserID        int64   `json:"user_id,omitempty"`
	RawMessage    string  `json:"raw_message,omitempty"`
	GroupID       int64   `json:"group_id,omitempty"`
	Time          int64   `json:"time"`
	MetaEventType string  `json:"meta_event_type,omitempty"`
	Sender        SendDer `json:"sender"`
	MessageID     int64   `json:"message_id"`
	MessageSeq    int64   `json:"message_seq"`
	MessageType   string  `json:"message_type"`       // "private" 或 "group"
	SubType       string  `json:"sub_type,omitempty"` // 私聊才有 "friend"
}

type SendDer struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nickname"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}

type NorMessageModel struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// PixivImage 结构体表示单个图片信息
type PixivImage struct {
	ID     uint     `gorm:"primaryKey;autoIncrement" json:"id"`      // 设为主键，自增
	PID    int      `gorm:"uniqueIndex" json:"pid"`                  // 唯一索引
	Title  string   `gorm:"type:varchar(255);not null" json:"title"` // 设定字符串长度
	Author string   `gorm:"type:varchar(255);not null" json:"author"`
	Tags   []string `gorm:"serializer:json" json:"tags"` // JSON 序列化存储

	URLs struct {
		Original string `json:"original"`
	} `gorm:"embedded" json:"urls"` // 使用 embedded 方式存储
}

type APIResponse struct {
	Error string       `json:"error"`
	Data  []PixivImage `json:"data"`
} //type Data struct {
//	interface{} `json:"data"`
//}
