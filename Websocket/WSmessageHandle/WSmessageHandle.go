package WSmessageHandle

import (
	"NepcatGoApiReq/ResHandle"
	"NepcatGoApiReq/Websocket"
)

var MessageHandlerFunc func(string) // 定义回调函数
// 处理普通 WebSocket 消息
func MessageHandler() {
	for msg := range Websocket.MessageChannel {
		//if MessageHandlerFunc != nil {
		ResHandle.HandleMessage(msg) // 触发回调，而不是直接调用 HandleDeepseekMessage
		//}
	}
}

var DeepSeekMessageHandlerFunc func(string) // 定义回调函数

func DeepSeekMessageHandler() {
	//for msg := range Websocket.DeepseekmessageChannel {
	//	//if DeepSeekMessageHandlerFunc != nil {
	//	//DeepSeekReqHandle.HandleDeepseekMessage(msg) // 触发回调，而不是直接调用 HandleDeepseekMessage
	//	//}
	//}
}
