package Websocket

import (
	"NepcatGoApiReq/ResHandle"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

// 消息处理通道
var messageChannel = make(chan string, 100) // 设定缓存大小，避免阻塞

func WebSocketInit() {
	serverURL := url.URL{
		Scheme:   "ws",             // WebSocket 协议
		Host:     "127.0.0.1:3001", // 服务器地址和端口
		Path:     "/",              // WebSocket 连接路径
		RawQuery: "access_token=",  // 这里可以填入你的 Token
	}

	// 建立 WebSocket 连接
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatalf("连接 WebSocket 失败: %v", err)
	}
	defer conn.Close()

	fmt.Println("✅ 成功连接到 WebSocket 服务器")

	// 捕获 Ctrl+C 退出
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 监听 WebSocket 消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("❌ 读取消息失败:", err)
				return
			}
			messages := string(message)
			fmt.Println(messages)
			fmt.Println("📩 收到消息:", string(message))
			// 发送到 channel 进行异步处理
			messageChannel <- string(message)
		}
	}()

	// 发送 WebSocket 消息
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"action":"get_login_info"}`))
	if err != nil {
		log.Println("❌ 发送消息失败:", err)
		return
	}
	fmt.Println("📤 已发送请求: 获取机器人登录信息")

	// 等待 Ctrl+C 退出
	<-interrupt
	fmt.Println("⏳ 关闭 WebSocket 连接...")
}

// 处理消息的 goroutine
func MessageHandler() {
	for msg := range messageChannel { // 持续监听 channel
		ResHandle.HandleMessage(msg) // 处理消息
	}
}
