package GameLogic

import (
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"fmt"
	"sort"
	"strings"
)

type GameManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *GameManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"用户注册":   userRegister,
		"获取宠物信息": getPetInfo,
		"等级查询":   levelQuery,
		"每日签到":   dailySignIn,
	}
	n.Handler = groupManagekeywordHandlers
}

// **获取按长度排序的关键词**
func (n *GameManageHandle) getSortedKeywords() []string {
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
func (n *GameManageHandle) HandleGameManageMessage(message MessageModel.Message) bool {
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

func userRegister(message MessageModel.Message) {
	fmt.Println(message)
	_, err := Tool.ExtractNumber(message.RawMessage)
	if err != nil {
		return
	}

}

func getPetInfo(message MessageModel.Message) {

}

func levelQuery(message MessageModel.Message) {

}

func dailySignIn(message MessageModel.Message) {

}
