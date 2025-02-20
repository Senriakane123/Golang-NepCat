package DSReqModel

// API 请求结构
type ChatRequest struct {
	Model    string      `json:"model"`
	Messages ChatMessage `json:"messages"`
	Stream   bool        `json:"stream"`
	MemoryID string      `json:"memory_id"`
}

// 消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// API 响应结构
type ChatResponseChunk struct {
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Choice struct {
	Delta Delta `json:"delta"`
}

type Delta struct {
	ReasoningContent *string `json:"reasoning_content,omitempty"`
	Content          *string `json:"content,omitempty"`
}

type Usage struct {
	TotalTokens int `json:"total_tokens"`
}

type RgsGroup struct {
	ID         int    `gorm:"primaryKey;autoIncrement;column:ID"`
	GroupID    int    `gorm:"type:int;not null;column:GroupID"`
	SeessionID string `gorm:"type:varchar(255);not null;column:SeessionID"` // 表中字段名为 skill_name
}

// 结构体: 创建 Memory 的 API 响应
type CreateMemoryResponse struct {
	MemoryID  string `json:"memoryId"`
	RequestID string `json:"requestId"`
}

// API 请求结构体
type RequestBody struct {
	Prompt    string `json:"prompt"`
	SessionID string `json:"session_id,omitempty"` // 可选的 session_id
}

// API 响应结构体
type ResponseBody struct {
	Output struct {
		Text      string `json:"text"`
		SessionID string `json:"session_id"`
	} `json:"output"`
	Usage struct {
		Models []struct {
			OutputTokens int    `json:"output_tokens"`
			ModelID      string `json:"model_id"`
			InputTokens  int    `json:"input_tokens"`
		} `json:"models"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
	Status    int    `json:"status_code"`
	Message   string `json:"message"`
}
