package DeepSeekReqHandle

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle/WebSocketControl"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"NepcatGoApiReq/Websocket"
	"context"
	"encoding/json"
	"fmt"
	"github.com/p9966/go-deepseek"
	"log"
	"sort"
	"strings"
)

var NowConnDSGroup int64

type DeepSeekManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *DeepSeekManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"接入deepseek": WebSocketControl.LoginInDeepSeek,
		"退出deepseek": WebSocketControl.LoginOutDeepSeek,
	}
	n.Handler = groupManagekeywordHandlers
}

// **获取按长度排序的关键词**
func (n *DeepSeekManageHandle) getSortedKeywords() []string {
	keys := make([]string, 0, len(n.Handler))
	for key := range n.Handler {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}

// 统一处理消息
func (n *DeepSeekManageHandle) HandleGroupManageMessage(message MessageModel.Message) bool {
	sortedKeywords := n.getSortedKeywords() // 获取按长度排序的关键词
	for _, keyword := range sortedKeywords {
		if strings.HasPrefix(message.RawMessage, keyword) || strings.Contains(message.RawMessage, keyword) {
			handler := n.Handler[keyword]
			NowConnDSGroup = message.GroupID
			handler(message)
			return true // 处理完一个就返回，避免重复触发
		}
	}
	return false
}

type DeepseekClient struct {
	Client deepseek.Client
}

func InitDeepSeekHandler() {
	//var DSClient DeepseekClient
	// 这里注册回调，避免循环导入
	Websocket.DeepSeekMessageHandlerFunc = HandleDeepseekMessage
}

var messageHistory []deepseek.OllamaChatMessage // 保存消息历史

//func HandleDeepseekMessage(message string) {
//	var msg MessageModel.Message
//	err := json.Unmarshal([]byte(message), &msg)
//	if err != nil {
//		fmt.Println("解析 JSON 失败:", err)
//		return
//	}
//	ResMsg, Boolen := Tool.ExtractDeepseekResMessage(msg.SelfID, msg.RawMessage)
//	if msg.PostType != "message" || msg.GroupID != NowConnDSGroup || !Boolen {
//		return
//	}
//
//	var DSHandle DeepSeekManageHandle
//	DSHandle.HandlerInit()
//	DSHandle.HandleGroupManageMessage(msg)
//
//	client := deepseek.Client{
//		BaseUrl: "http://localhost:11434",
//	}
//
//	request := deepseek.OllamaChatRequest{
//		Model: "deepseek-r1:8b",
//		Messages: []deepseek.OllamaChatMessage{
//			{
//				Role:    "user",
//				Content: ResMsg,
//			},
//		},
//	}
//
//	response, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
//	if err != nil {
//		log.Fatalf("Error: %v", err)
//	}
//
//	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
//		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(msg.GroupID, "[CQ:at,qq="+Tool.Int64toString(msg.Sender.UserID)+"]"+response.Message.Content))
//	}
//	fmt.Println(response.Message.Content)
//}

func HandleDeepseekMessage(message string) {
	var msg MessageModel.Message
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}

	ResMsg, Boolen := Tool.ExtractDeepseekResMessage(msg.SelfID, msg.RawMessage)
	if msg.PostType != "message" || msg.GroupID != NowConnDSGroup || !Boolen {
		return
	}

	var DSHandle DeepSeekManageHandle
	DSHandle.HandlerInit()
	DSHandle.HandleGroupManageMessage(msg)

	// 更新消息历史，保存当前消息
	messageHistory = append(messageHistory, deepseek.OllamaChatMessage{
		Role:    "user",
		Content: ResMsg,
	})

	client := deepseek.Client{
		BaseUrl: "http://localhost:11434",
	}

	// 使用历史消息来构建请求
	request := deepseek.OllamaChatRequest{
		Model:    "deepseek-r1:8b",
		Messages: messageHistory, // 传递整个历史消息
	}

	// 发送请求并获取响应
	response, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 将模型的回复加入历史记录
	messageHistory = append(messageHistory, deepseek.OllamaChatMessage{
		Role:    "assistant",
		Content: response.Message.Content,
	})

	// 发送响应到群聊
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(msg.GroupID, "[CQ:at,qq="+Tool.Int64toString(msg.Sender.UserID)+"]"+response.Message.Content))
	}

	fmt.Println(response.Message.Content)
}
