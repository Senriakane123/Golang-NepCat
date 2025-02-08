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
	NickName string `json:"nick_name"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}

type NorMessageModel struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// PixivImage 结构体表示单个图片信息
type PixivImage struct {
	PID    int      `json:"pid"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Tags   []string `json:"tags"`
	URLs   struct {
		Original string `json:"original"`
	} `json:"urls"`
}

type APIResponse struct {
	Error string       `json:"error"`
	Data  []PixivImage `json:"data"`
} //type Data struct {
//	interface{} `json:"data"`
//}
