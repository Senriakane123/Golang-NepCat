package DeepSeekReqHandle

import (
	"NepcatGoApiReq/Database/DBControlApi"
	"NepcatGoApiReq/Database/DBModel"
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/ReqApiConst"
	"io"

	"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle/DSReqModel"
	//"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle/MemoryIDCtl"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	//"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/p9966/go-deepseek"
	"io/ioutil"
	"net/http"
	//"strings"
)

//const apiURL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
//const apiKey = "sk-xxx"                 // 这里填入你的 DashScope API Key
//const historyFile = "chat_history.json" // 本地存储对话历史

var NowConnDSGroup int64

type DeepSeekManageHandle struct {
	//Handler map[string]func(message MessageModel.Message)
	apiURL      string
	apiKey      string
	historyFile string
	Handler     func(message MessageModel.Message)
	maxHistory  int64
}

func (n *DeepSeekManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	//var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
	//	"接入deepseek": WebSocketControl.LoginInDeepSeek,
	//	"退出deepseek": WebSocketControl.LoginOutDeepSeek,
	//}
	//n.Handler = groupManagekeywordHandlers

	//n.Handler = WebSocketControl.LoginInDeepSeek

	n.Handler = n.HandleCloudDeepseekMessage
	//n.maxHistory = 300
	//n.apiURL = "https://dashscope.aliyuncs.com/api/v1/apps/99d007044b5849e08b145ee2bfd6f174/completion"

}

// **获取按长度排序的关键词**
//func (n *DeepSeekManageHandle) getSortedKeywords() []string {
//	keys := make([]string, 0, len(n.Handler))
//	for key := range n.Handler {
//		keys = append(keys, key)
//	}
//	sort.Slice(keys, func(i, j int) bool {
//		return len(keys[i]) > len(keys[j])
//	})
//	return keys
//}

// 统一处理消息
func (n *DeepSeekManageHandle) HandleGroupManageMessage(message MessageModel.Message) bool {
	//sortedKeywords := n.getSortedKeywords() // 获取按长度排序的关键词
	//for _, keyword := range sortedKeywords {
	//	if strings.HasPrefix(message.RawMessage, keyword) || strings.Contains(message.RawMessage, keyword) {
	//		handler := n.Handler[keyword]
	//		//第一次进行ds接入会把群号信息存入NowConnDSGroup
	//		NowConnDSGroup = message.GroupID
	//		handler(message)
	//		return true // 处理完一个就返回，避免重复触发
	//	}
	//}
	//n.historyFile = Tool.Int64toString(message.GroupID) + "_chat_history.json"
	var AdminUser []DBModel.AdminUser
	//var sessionID string
	_, err := DBControlApi.Db.Where("adminuser", &AdminUser, "QQNum = ?", message.Sender.UserID)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "验证错误，或用户不具有管理权限"))
		}
	} else {
		if len(AdminUser) == 0 {
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
				handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "用户不具有管理权限"))
			}
		}
		for _, v := range AdminUser {
			if v.QQNum == int(message.Sender.UserID) {
				n.Handler(message)
				return true
			}
		}
	}
	return true
}

type DeepseekClient struct {
	Client deepseek.Client
}

//func InitDeepSeekHandler() {
//	//var DSClient DeepseekClient
//	// 这里注册回调，避免循环导入
//	Websocket.DeepSeekMessageHandlerFunc = HandleDeepseekMessage
//}

var messageHistory []deepseek.OllamaChatMessage // 保存消息历史

//func HandleDeepseekMessage(msg MessageModel.Message) {
//	//var msg MessageModel.Message
//	//err := json.Unmarshal([]byte(message), &msg)
//	//if err != nil {
//	//	fmt.Println("解析 JSON 失败:", err)
//	//	return
//	//}
//	//判断消息类型和消息的
//	ResMsg, Boolen := Tool.ExtractDeepseekResMessage(msg.SelfID, msg.RawMessage)
//	if msg.PostType != "message" || msg.GroupID != NowConnDSGroup || !Boolen {
//		return
//	}
//	//这里疑似冗余代码
//	//////////////////////////////////////////////////
//	//////////////////////////////////////////////////
//	//////////////////////////////////////////////////
//	var DSHandle DeepSeekManageHandle
//	DSHandle.HandlerInit()
//	DSHandle.HandleGroupManageMessage(msg)
//	//////////////////////////////////////////////////
//	//////////////////////////////////////////////////
//	//////////////////////////////////////////////////
//
//	// 更新消息历史，保存当前消息
//	messageHistory = append(messageHistory, deepseek.OllamaChatMessage{
//		Role:    "user",
//		Content: ResMsg,
//	})
//
//	client := deepseek.Client{
//		BaseUrl: "http://localhost:11434",
//	}
//
//	// 使用历史消息来构建请求
//	request := deepseek.OllamaChatRequest{
//		Model:    "deepseek-r1:8b",
//		Messages: messageHistory, // 传递整个历史消息
//	}
//
//	// 发送请求并获取响应
//	response, err := client.CreateOllamaChatCompletion(context.TODO(), &request)
//	if err != nil {
//		log.Fatalf("Error: %v", err)
//	}
//
//	// 将模型的回复加入历史记录
//	messageHistory = append(messageHistory, deepseek.OllamaChatMessage{
//		Role:    "assistant",
//		Content: response.Message.Content,
//	})
//
//	// 发送响应到群聊
//	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
//		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(msg.GroupID, "[CQ:at,qq="+Tool.Int64toString(msg.Sender.UserID)+"]"+response.Message.Content))
//	}
//
//	fmt.Println(response.Message.Content)
//}

func (n *DeepSeekManageHandle) HandleCloudDeepseekMessage(msg MessageModel.Message) {

	var rgsgroup []DSReqModel.RgsGroup
	var sessionID string
	_, err := DBControlApi.Db.Where("rgsgroup", &rgsgroup, "GroupID = ?", msg.GroupID)
	if err != nil {
		sessionID = ""
	} else {
		if len(rgsgroup) == 0 {
			sessionID = ""
		} else {
			sessionID = rgsgroup[0].SeessionID
		}

	}

	// 读取环境变量中的 API Key
	//apiKey := os.Getenv("DASHSCOPE_API_KEY") // 确保你的环境变量已设置
	//if apiKey == "" {
	//	fmt.Println("Error: DASHSCOPE_API_KEY not set")
	//	return
	//}

	apiKey := "sk-efd428b2a41f42668dc8579cc6536281"
	appID := "99d007044b5849e08b145ee2bfd6f174" // 替换为你的 APP_ID
	url := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/apps/%s/completion", appID)

	// 创建请求体
	requestBody := map[string]interface{}{
		"input": map[string]string{
			"prompt":     Tool.Int64toString(msg.Sender.UserID) + "-" + msg.Sender.NickName + "：" + Tool.ExtractMessageForRob(msg.RawMessage),
			"session_id": sessionID, // 替换为实际上一轮对话的session_id
		},
		"parameters": map[string]interface{}{},
		"debug":      map[string]interface{}{},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		return
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	// 处理响应
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request successful:")
		fmt.Println(string(body))
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		fmt.Println(string(body))
	}

	// 解析 JSON 响应
	var responseBody DSReqModel.ResponseBody
	err = json.Unmarshal(body, &responseBody)
	//fmt.Println(string(body))
	//fmt.Println(responseBody)
	if err != nil {
		return
	}

	if len(rgsgroup) == 0 {
		//var newGroupHandle DSReqModel.RgsGroup
		newGroupHandle := DSReqModel.RgsGroup{
			GroupID:    int(msg.GroupID),
			SeessionID: responseBody.Output.SessionID,
		}
		fmt.Println(responseBody.Output.SessionID)
		err = DBControlApi.Db.Create(&newGroupHandle, "rgsgroup")
		if err != nil {
			fmt.Printf("Failed to update rgsgroup: %v\n", err)
		}
	}
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(msg.GroupID, responseBody.Output.Text))
	}
	return

}

// 读取历史对话
func (n *DeepSeekManageHandle) loadHistory() []DSReqModel.ChatMessage {
	data, err := ioutil.ReadFile(n.historyFile)
	if err != nil {
		// 文件不存在，返回空历史
		return []DSReqModel.ChatMessage{}
	}
	var messages []DSReqModel.ChatMessage
	if err = json.Unmarshal(data, &messages); err != nil {
		fmt.Println("解析历史记录失败:", err)
		return []DSReqModel.ChatMessage{}
	}
	return messages
}

// 存储历史对话
func (n *DeepSeekManageHandle) saveHistory(messages []DSReqModel.ChatMessage) {
	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		fmt.Println("存储历史记录失败:", err)
		return
	}
	_ = ioutil.WriteFile(n.historyFile, data, 0644)
}

func (n *DeepSeekManageHandle) trimHistory(messages []DSReqModel.ChatMessage) []DSReqModel.ChatMessage {
	if len(messages) > int(n.maxHistory) {
		return messages[len(messages)-int(n.maxHistory):] // 只保留最后 maxHistory 条
	}
	return messages
}
