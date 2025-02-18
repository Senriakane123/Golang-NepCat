package WebSocketControl

import (
	"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle/NepcatWebSocketServerCtr"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/NepCatSysHttpReq"
	"NepcatGoApiReq/Websocket"
	"encoding/json"
	"fmt"
	"strconv"
)

func LoginInDeepSeek(message MessageModel.Message) {
	var WsHttphandle NepCatSysHttpReq.NepCatHttpReq
	WsHttphandle.BaseUrlAndPasswordinit("cx030115", "http://127.0.0.1:6099")
	err := WsHttphandle.Login()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	netcofig, err := WsHttphandle.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ReqNetconfig, NewWsServerInfo := NepcatWebSocketServerCtr.AddWsServer(netcofig)

	jsonData, err := json.Marshal(ReqNetconfig)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}

	WsHttphandle.Config.Config = string(jsonData)

	err = WsHttphandle.SetConfig()
	if err != nil {
		fmt.Println("创建新ws服务器失败", err)
		return
	}

	go Websocket.DeepSeekMessageHandler()
	Websocket.CloseWebSocket() // 关闭默认 WebSocket
	//CloseDeepSeekWebSocket() // 关闭之前的 DeepSeek 连接

	Websocket.WebSocketInitForDeepSeek(strconv.Itoa(NewWsServerInfo.Port)) // 先关闭普通 WebSocket，再连接 DeepSeek
}

func LoginOutDeepSeek(message MessageModel.Message) {
	fmt.Println("🔴 退出 DeepSeek WebSocket，恢复默认 WebSocket...")

	// 关闭 DeepSeek WebSocket
	Websocket.CloseDeepSeekWebSocket()

	go Websocket.MessageHandler()
	// 重新启动默认 WebSocket
	Websocket.WebSocketInit()
}
