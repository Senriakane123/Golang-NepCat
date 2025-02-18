package NepcatWebSocketServerCtr

import (
	"NepcatGoApiReq/NepCatSysHttpReq/NepcatReqModel"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

//func AddWsServer(wscofnig string) NepcatReqModel.Config {
//	var configResponse NepcatReqModel.ConfigResponse
//
//	// 解析 JSON
//	err := json.Unmarshal([]byte(wscofnig), &configResponse)
//	if err != nil {
//		fmt.Println("JSON 解析失败:", err)
//		return
//	}
//
//	for _, v := range configResponse.Data.Network.WebsocketServers {
//
//	}
//
//	// 输出解析后的数据
//	fmt.Printf("解析后的结构体: %+v\n", configResponse)
//
//	return Netconfig
//}

func AddWsServer(wsconfig string) (NepcatReqModel.Config, NepcatReqModel.WebsocketServer) {
	var configResponse NepcatReqModel.ConfigResponse
	// 解析 JSON
	err := json.Unmarshal([]byte(wsconfig), &configResponse)
	if err != nil {
		fmt.Println("JSON 解析失败:", err)
		return NepcatReqModel.Config{}, NepcatReqModel.WebsocketServer{}
	}

	// 获取现有端口列表
	existingPorts := make(map[int]bool)
	for _, v := range configResponse.Data.Network.WebsocketServers {
		existingPorts[v.Port] = true
	}

	// 生成一个不重复的随机端口
	rand.Seed(time.Now().UnixNano())
	var newPort int
	for {
		newPort = rand.Intn(10000) + 2000 // 生成 2000-11999 之间的端口
		if !existingPorts[newPort] {
			break
		}
	}

	// 创建新的 WebSocket 服务器配置
	newWsServer := NepcatReqModel.WebsocketServer{
		Enable:               true,
		Name:                 fmt.Sprintf("DSRandomWSServer-%d", newPort),
		Host:                 "0.0.0.0",
		Port:                 newPort,
		ReportSelfMessage:    false,
		EnableForcePushEvent: true,
		MessagePostFormat:    "array",
		Token:                "",
		Debug:                true,
		HeartInterval:        30000,
	}

	// 添加到 WebSocket 服务器列表
	configResponse.Data.Network.WebsocketServers = append(configResponse.Data.Network.WebsocketServers, newWsServer)

	// 输出解析后的结构体
	fmt.Printf("添加新的 WebSocket 服务器后: %+v\n", configResponse.Data.Network.WebsocketServers)

	return configResponse.Data, newWsServer
}
