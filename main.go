package main

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/Websocket"
)

func main() {
	HTTPReq.InitAllApis()
	//DBControlApi
	//DBControlApi.InitDatabase()
	go Websocket.MessageHandler()

	Websocket.WebSocketInit()

}
