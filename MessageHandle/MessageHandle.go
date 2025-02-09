package MessageHandle

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"fmt"
	"strings"
)

func GroupMessageHandle(message MessageModel.Message) {

	_, QQNumberList := AtQQNumber(message.RawMessage)
	//先判断整体是否包含@bot的字符串
	if strings.Contains(message.RawMessage, "[CQ:at,qq=3666859102]") {
		fmt.Println(message.RawMessage)
		//解析是否需要服务
		ServerMach := ParseServiceCommand(message.RawMessage)
		if len(ServerMach) == 0 {
			//处理各种禁言和群权限管理
			//if HandleGroupManageMessage(message) {
			//	return
			//}
			//群管理功能包
			var GroupMessageAPi GroupManageHandle
			GroupMessageAPi.HandlerInit()
			if GroupMessageAPi.HandleGroupManageMessage(message) {
				return
			}

			var RandomPic PicManage
			RandomPic.HandlerInit()
			if RandomPic.HandlePicManageMessage(message) {
				return
			}

			//获取初级功能菜单
			for _, QQNumber := range QQNumberList {
				if message.SelfID == StringToInt64(QQNumber) {
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, BuildReplyMessage(message.UserID, MessageModel.GetServerList()), message.GroupID, MessageModel.NormalRespMessage(message.GroupID, BuildReplyMessage(message.UserID, MessageModel.GetServerList())))
					}
				}

			}
		} else {
			//判断初级服务菜单
			if ServerMach[2] == "" {
				switch ServerMach[1] {
				case "1":
				case "2":
				case "3":
				case "4":
				case "5":
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList5()), message.GroupID, MessageModel.NormalRespMessage(message.GroupID, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList5())))
					}
				case "6":
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList6()), message.GroupID, MessageModel.NormalRespMessage(message.GroupID, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList6())))
					}
				case "7":
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList6()), message.GroupID, MessageModel.NormalRespMessage(message.GroupID, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList6())))
					}
				}

			}
		}

	}

}

func ErrMessageHandle(Repmessage string, message MessageModel.Message) {
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList6()), message.GroupID, MessageModel.NormalRespMessage(message.GroupID, BuildReplyMessage(message.UserID, MessageModel.GetChildServerList6())))
	}
}
