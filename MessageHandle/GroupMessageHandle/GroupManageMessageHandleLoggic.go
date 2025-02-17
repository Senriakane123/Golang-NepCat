package GroupMessageHandle

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"sort"
	"strings"
)

type GroupManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *GroupManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"禁言":     n.handleBan,
		"解除全体禁言": n.handleUnbanGroup,
		"全体禁言":   n.handleGroupBan,
		"踢人":     n.KickSomeBody,
	}
	n.Handler = groupManagekeywordHandlers
}

// **获取按长度排序的关键词**
func (n *GroupManageHandle) getSortedKeywords() []string {
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
func (n *GroupManageHandle) HandleGroupManageMessage(message MessageModel.Message) bool {
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

// 处理全体禁言
func (n *GroupManageHandle) handleGroupBan(message MessageModel.Message) {
	if message.Sender.Role == "owner" || message.Sender.Role == "admin" {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SET_GROUP_WHOLE_BAN]; exists {
			handler(ReqApiConst.SET_GROUP_WHOLE_BAN, MessageModel.GroupBanRespMessage(message.GroupID, true))
		}
	} else {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]"+"请求格式错误或用户不具备管理员权限"))
		}
		return
	}

}

// 处理解除全体禁言
func (n *GroupManageHandle) handleUnbanGroup(message MessageModel.Message) {
	if message.Sender.Role == "owner" || message.Sender.Role == "admin" {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SET_GROUP_WHOLE_BAN]; exists {
			handler(ReqApiConst.SET_GROUP_WHOLE_BAN, MessageModel.GroupBanRespMessage(message.GroupID, false))
		}
	} else {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]"+"请求格式错误或用户不具备管理员权限"))
		}
		return
	}

}

func (n *GroupManageHandle) KickSomeBody(message MessageModel.Message) {
	if message.Sender.Role == "owner" || message.Sender.Role == "admin" {
		_, QQNumberList := Tool.AtQQNumber(message.RawMessage)
		for _, v := range QQNumberList {
			if Tool.StringToInt64(v) == message.SelfID {
				continue
			}
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SET_GROUP_KICK]; exists {

				handler(ReqApiConst.SET_GROUP_KICK, MessageModel.KickRespMessage(message.GroupID, Tool.StringToInt64(v), false))
			}
			return
		}
	} else {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]"+"请求格式错误或用户不具备管理员权限"))
		}
		return
	}

}

//

// 处理单独禁言
func (n *GroupManageHandle) handleBan(message MessageModel.Message) {
	_, QQNumberList := Tool.AtQQNumber(message.RawMessage)
	BanReqFormatBoolen, _, Time := Tool.CheckBanFormat(message.RawMessage)
	if (message.Sender.Role == "owner" || message.Sender.Role == "admin") && BanReqFormatBoolen {
		for _, v := range QQNumberList {
			if Tool.StringToInt64(v) == message.SelfID {
				continue
			}
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SET_GROUP_BAN]; exists {
				//for _, n := range BanQQNumberList {
				//
				//}
				handler(ReqApiConst.SET_GROUP_BAN, MessageModel.BanRespMessage(message.GroupID, Tool.StringToInt64(v), int64(Time)))
			}
			return
		}
	} else {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]"+"请求格式错误或用户不具备管理员权限"))
		}
		return
	}
}
