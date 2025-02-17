package WebSocketControl

import (
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/Websocket"
	"fmt"
)

func LoginInDeepSeek(message MessageModel.Message) {
	
	go Websocket.DeepSeekMessageHandler()
	Websocket.CloseWebSocket() // 关闭默认 WebSocket
	//CloseDeepSeekWebSocket() // 关闭之前的 DeepSeek 连接
	Websocket.WebSocketInitForDeepSeek() // 先关闭普通 WebSocket，再连接 DeepSeek
}

func LoginOutDeepSeek(message MessageModel.Message) {
	fmt.Println("🔴 退出 DeepSeek WebSocket，恢复默认 WebSocket...")

	// 关闭 DeepSeek WebSocket
	Websocket.CloseDeepSeekWebSocket()

	go Websocket.MessageHandler()
	// 重新启动默认 WebSocket
	Websocket.WebSocketInit()
}
