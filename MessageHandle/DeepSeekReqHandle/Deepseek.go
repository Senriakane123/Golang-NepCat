package DeepSeekReqHandle

import (
	"NepcatGoApiReq/MessageModel"
	"sort"
	"strings"
)

type DeepSeekManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *DeepSeekManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"接入deepseek": loginInDeepSeek,
		"解除deepseek": loginOutDeepSeek,
	}
	n.Handler = groupManagekeywordHandlers
}

// **获取按长度排序的关键词**
func (n *DeepSeekManageHandle) getSortedKeywords() []string {
	keys := make([]string, 0, len(n.Handler))
	for key := range n.Handler {
		keys = append(keys, key)
	}
	// 按字符串长度从长到短排序，保证 "解除全体禁言" 先匹配，而不是 "禁言"
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}

// 统一处理消息
func (n *DeepSeekManageHandle) HandleGroupManageMessage(message MessageModel.Message) bool {
	sortedKeywords := n.getSortedKeywords() // 获取按长度排序的关键词
	for _, keyword := range sortedKeywords {
		if strings.HasPrefix(message.RawMessage, keyword) || strings.Contains(message.RawMessage, keyword) {
			handler := n.Handler[keyword]
			handler(message)
			return true // 处理完一个就返回，避免重复触发
		}
	}
	return false
}

func loginInDeepSeek(message MessageModel.Message) {

}

func loginOutDeepSeek(message MessageModel.Message) {

}
