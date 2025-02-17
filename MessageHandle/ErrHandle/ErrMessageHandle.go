package ErrHandle

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
)

func ErrMessageHandle(Repmessage string, message MessageModel.Message) {
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Repmessage))
	}
}
