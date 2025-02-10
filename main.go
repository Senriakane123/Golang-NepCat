package main

import (
	"NepcatGoApiReq/Database/DBControlApi"
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/Websocket"
)

func main() {
	HTTPReq.InitAllApis()
	//DBControlApi
	DBControlApi.InitDatabase()
	go Websocket.MessageHandler()

	Websocket.WebSocketInit()

}
