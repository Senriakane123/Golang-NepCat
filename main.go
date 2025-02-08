package main

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/Websocket"
)

func main() {
	HTTPReq.InitAllApis()
	go Websocket.MessageHandler()

	Websocket.WebSocketInit()

}
