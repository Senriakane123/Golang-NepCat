package Websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
)

// 消息处理通道
var MessageChannel = make(chan string, 100)
var DeepseekmessageChannel = make(chan string, 100)

//type WSMgrHandle struct {
//	conn *websocket.Conn
//	//MessageChannel  make(chan string, 100)
//}
//
//type DSWSMgrHandle struct {
//	deepseekConn *websocket.Conn
//	mu           sync.Mutex // 互斥锁，避免并发问题
//}
//
//var WSHandle WSMgrHandle
//var DSWSHandle DSWSMgrHandle

// WebSocket 连接实例（用于管理多个 WebSocket 连接）
var conn *websocket.Conn
var deepseekConn *websocket.Conn
var mu sync.Mutex // 互斥锁，避免并发问题

// 关闭当前 WebSocket 连接
func CloseWebSocket() {
	mu.Lock()
	defer mu.Unlock()
	if conn != nil {
		_ = conn.Close()
		fmt.Println("🔴 WebSocket 连接已关闭")
		conn = nil
	}
}

// 关闭 DeepSeek WebSocket 连接
func CloseDeepSeekWebSocket() {
	mu.Lock()
	defer mu.Unlock()
	if deepseekConn != nil {
		_ = deepseekConn.Close()
		fmt.Println("🔴 DeepSeek WebSocket 连接已关闭")
		deepseekConn = nil
	}
}

// 初始化 WebSocket（默认连接）
func WebSocketInit() {

	serverURL := url.URL{
		Scheme:   "ws",
		Host:     "127.0.0.1:3001",
		Path:     "/",
		RawQuery: "access_token=",
	}

	var err error
	conn, _, err = websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatalf("❌ 连接 WebSocket 失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 WebSocket 服务器")

	// 捕获 Ctrl+C 退出
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("❌ 读取消息失败:", err)
				return
			}
			fmt.Println("📩 收到消息:", string(message))
			MessageChannel <- string(message)
			//go MessageHandler()
		}
	}()

	// 等待 Ctrl+C 退出
	<-interrupt
	fmt.Println("⏳ 关闭 WebSocket 连接...")
}

// 初始化 DeepSeek WebSocket
func WebSocketInitForDeepSeek() {

	serverURL := url.URL{
		Scheme:   "ws",
		Host:     "127.0.0.1:3002",
		Path:     "/",
		RawQuery: "access_token=",
	}

	var err error
	deepseekConn, _, err = websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatalf("❌ 连接 DeepSeek WebSocket 失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 DeepSeek WebSocket 服务器")

	// 捕获 Ctrl+C 退出
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		//var nowConnGroup string
		for {
			_, message, err := deepseekConn.ReadMessage()
			if err != nil {
				log.Println("❌ 读取 DeepSeek 消息失败:", err)
				return
			}
			fmt.Println("📩 收到 DeepSeek 消息:", string(message))
			DeepseekmessageChannel <- string(message)
			//go DeepSeekMessageHandler()
		}
	}()

	// 等待 Ctrl+C 退出
	<-interrupt
	fmt.Println("⏳ 关闭 WebSocket 连接...")
}

var MessageHandlerFunc func(string) // 定义回调函数
// 处理普通 WebSocket 消息
func MessageHandler() {
	for msg := range MessageChannel {
		if MessageHandlerFunc != nil {
			MessageHandlerFunc(msg) // 触发回调，而不是直接调用 HandleDeepseekMessage
		}
		//return
	}

}

var DeepSeekMessageHandlerFunc func(string) // 定义回调函数

func DeepSeekMessageHandler() {
	for msg := range DeepseekmessageChannel {
		if DeepSeekMessageHandlerFunc != nil {
			DeepSeekMessageHandlerFunc(msg) // 触发回调，而不是直接调用 HandleDeepseekMessage
		}
	}
}
