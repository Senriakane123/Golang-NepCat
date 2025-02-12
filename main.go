package main

import (
	"NepcatGoApiReq/Database/DBControlApi"
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle"
	"NepcatGoApiReq/ResHandle"
	"NepcatGoApiReq/Websocket"
)

func main() {
	ResHandle.InitDeepSeekHandler()         //初始化normal消息处理
	DeepSeekReqHandle.InitDeepSeekHandler() // 初始化 DeepSeek 消息处理
	HTTPReq.InitAllApis()
	//DBControlApi
	DBControlApi.InitDatabase() //数据库注册
	go Websocket.MessageHandler()

	Websocket.WebSocketInit()

}
