package MessageHandle

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/GroupMessageHandle"
	"NepcatGoApiReq/MessageHandle/R18PicManage"
	"NepcatGoApiReq/MessageHandle/RPGGameHandle/GameLogic"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"fmt"
	"strings"
)

func MessageHandle(message MessageModel.Message) {

	_, QQNumberList := Tool.AtQQNumber(message.RawMessage)
	//先判断整体是否包含@bot的字符串
	if strings.Contains(message.RawMessage, "[CQ:at,qq=3666859102]") {
		fmt.Println(message.RawMessage)
		//解析是否需要服务
		ServerMach := Tool.ParseServiceCommand(message.RawMessage)
		if len(ServerMach) == 0 {

			//游戏功能包
			var GameMessageApi GameLogic.GameManageHandle
			GameMessageApi.HandlerInit()
			if GameMessageApi.HandleGameManageMessage(message) {
				return
			}

			//群管理功能包
			var GroupMessageAPi GroupMessageHandle.GroupManageHandle
			GroupMessageAPi.HandlerInit()
			if GroupMessageAPi.HandleGroupManageMessage(message) {
				return
			}

			//涩图功能包
			var RandomPic R18PicManage.PicManage
			RandomPic.HandlerInit()
			if RandomPic.HandlePicManageMessage(message) {
				return
			}

			//获取初级功能菜单
			for _, QQNumber := range QQNumberList {
				if message.SelfID == Tool.StringToInt64(QQNumber) {
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(MessageModel.GetServerList())))
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
						handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(MessageModel.GetChildServerList5())))
					}
				case "6":
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(MessageModel.GetChildServerList6())))
					}
				case "7":
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(MessageModel.GetChildServerList6())))
					}
				case "8":
					if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
						handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(MessageModel.GetChildServerList8())))
					}

				}

			}
		}

	}

}

func ErrMessageHandle(Repmessage string, message MessageModel.Message) {
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(MessageModel.GetChildServerList6())))
	}
}
