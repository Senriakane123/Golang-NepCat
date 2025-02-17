package RobInfoManage

import (
	"NepcatGoApiReq/Database/DBControlApi"
	"NepcatGoApiReq/Database/DBModel"
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"fmt"
	"sort"
	"strings"
)

type GroupManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *GroupManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"修改头像": n.ChangeRobPic,
		//"解除全体禁言": n.handleUnbanGroup,
		//"全体禁言":   n.handleGroupBan,
		//"踢人":     n.KickSomeBody,
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
func (n *GroupManageHandle) HandleRobInfoMessage(message MessageModel.Message) bool {
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

func (n *GroupManageHandle) ChangeRobPic(message MessageModel.Message) {
	var AdminUserlist []DBModel.AdminUser
	_, err := DBControlApi.Db.Where("adminuser", &AdminUserlist, "QQNum = ?", message.Sender.UserID)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]"+"查询管理数据表失败"))
		}
	} else {
		if len(AdminUserlist) > 0 {
			fmt.Println(message.RawMessage)
			link := Tool.ExtractLinks(message.RawMessage)
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SET_QQ_AVATAR]; exists {
				handler(ReqApiConst.SET_QQ_AVATAR, MessageModel.ChangeAvatarMessage(link[0]))
			}
		} else {
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
				handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]"+"用户不具有管理员权限"))
			}
		}

	}

}
