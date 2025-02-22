package DeepSeekReqHandle

import (
	"NepcatGoApiReq/Database/DBControlApi"
	"NepcatGoApiReq/Database/DBModel"
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle/DSReqModel"
	"NepcatGoApiReq/ReqApiConst"
	"crypto/rand"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"io"
	"log"
	"regexp"

	//"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle/MemoryIDCtl"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	//"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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
	Handler     map[string]func(message MessageModel.Message)
	maxHistory  int64
}

func (n *DeepSeekManageHandle) HandlerInit() {
	// 确保 map 被初始化
	n.Handler = make(map[string]func(message MessageModel.Message))

	n.Handler["CloudDS"] = n.HandleCloudDeepseekMessage
	n.Handler["LocalDS"] = n.HandleLocalDeepseekMessage
}

// 统一处理消息
func (n *DeepSeekManageHandle) HandleGroupManageMessage(message MessageModel.Message) bool {

	var AdminUser []DBModel.AdminUser
	//var sessionID string
	_, err := DBControlApi.Db.Where("adminuser", &AdminUser, "QQNum = ?", message.Sender.UserID)
	if err != nil {
		if message.MessageType == "private" {
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_PRIVATE_MSG]; exists {
				handler(ReqApiConst.SEND_PRIVATE_MSG, MessageModel.NormalPrivateRespMessage(message.Sender.UserID, "验证错误，或用户不具有管理权限,请请求最高管理员735439479获取管理权限"))
			}
			return false
		}
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "验证错误，或用户不具有管理权限,请请求最高管理员735439479获取管理权限"))
		}
		return false
	} else {
		if len(AdminUser) == 0 {
			//n.Handler = n.HandleLocalDeepseekMessage
			n.Handler["LocalDS"](message)
			return true
			//if message.MessageType == "private" {
			//	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_PRIVATE_MSG]; exists {
			//		handler(ReqApiConst.SEND_PRIVATE_MSG, MessageModel.NormalPrivateRespMessage(message.Sender.UserID, "用户不具有管理权限,请请求最高管理员735439479获取管理权限"))
			//	}
			//	return false
			//}
			//if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			//	handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "用户不具有管理权限,请请求最高管理员735439479获取管理权限"))
			//}
			//return false
		}
		for _, v := range AdminUser {
			if v.QQNum == int(message.Sender.UserID) {
				n.Handler["CloudDS"](message)
				return true
			}
		}
	}
	return true
}

func (n *DeepSeekManageHandle) HandleLocalDeepseekMessage(msg MessageModel.Message) {
	//1.判断私人还是群聊
	//2.更具连接类型，先查询对应数据库，查看是否第一次连接
	//3.第一次连接随机生成sessionID，存入数据库，不是第一次则会从数据库拿去sessionID
	//4.更具不同连接类型解析不同发送给rob的消息
	//5.组件发送rob请求
	//6.更具不同连接类型返回内容
	var rgsgroup []DSReqModel.RgsGroup
	var rgsprivate []DSReqModel.RgsPrivate
	var sessionID string
	var SendMsgtype string
	var sendMsgID int64
	var RobReqMsg string

	if msg.MessageType == "private" {
		_, err := DBControlApi.Db.Where("localrgsprivate", &rgsprivate, "QQID = ?", msg.Sender.UserID)
		if err != nil {
			sessionID = GenerateSessionID()
		} else {
			if len(rgsprivate) == 0 {
				sessionID = GenerateSessionID()
			} else {
				sessionID = rgsprivate[0].SessionID
			}

		}
		//设置发送的群或者人
		sendMsgID = msg.Sender.UserID
		//设置该调用的类型
		SendMsgtype = ReqApiConst.SEND_PRIVATE_MSG
		//定制送个rob消息类型
		RobReqMsg = msg.RawMessage
	} else if msg.MessageType == "group" {
		_, err := DBControlApi.Db.Where("localrgsgroup", &rgsgroup, "GroupID = ?", msg.GroupID)
		if err != nil {
			sessionID = GenerateSessionID()
		} else {
			if len(rgsgroup) == 0 {
				sessionID = GenerateSessionID()
			} else {
				sessionID = rgsgroup[0].SeessionID
			}

		}

		sendMsgID = msg.GroupID
		SendMsgtype = ReqApiConst.SEND_GROUP_MSG
		RobReqMsg = Tool.ExtractMessageForRob(msg.RawMessage)
	}

	client := resty.New()
	apiKey := "NH8TBX3-DXSM2Z0-M1VT66A-XC3QD8V"
	workspaceSlug := "forqqrobot" // **注意这里要用小写**
	url := fmt.Sprintf("http://localhost:54079/api/v1/workspace/%s/chat", workspaceSlug)
	//RobReqMsg = Tool.ExtractMessageForRob(msg.RawMessage)
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"message":   RobReqMsg,
			"mode":      "chat", // **必须是 "query" 或 "chat"，不能用 "query | chat"**
			"sessionId": sessionID,
		}).
		Post(url)

	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	// 解析 JSON 响应
	var responseBody DSReqModel.LocalResponse
	err = json.Unmarshal([]byte(resp.String()), &responseBody)
	//fmt.Println(string(body))
	//fmt.Println(responseBody)
	if err != nil {
		return
	}
	fmt.Println(responseBody.TextResponse)
	// 使用正则表达式去除 <think> 标签中的内容
	re := regexp.MustCompile(`<think>.*?</think>`)
	Resp := re.ReplaceAllString(responseBody.TextResponse, "")

	//fmt.Println("响应:", resp.String())

	if len(rgsgroup) == 0 && msg.MessageType == "group" {
		//var newGroupHandle DSReqModel.RgsGroup
		newGroupHandle := DSReqModel.RgsGroup{
			GroupID:    int(msg.GroupID),
			SeessionID: sessionID,
		}
		//fmt.Println(responseBody.Output.SessionID)
		err = DBControlApi.Db.Create(&newGroupHandle, "localrgsgroup")
		if err != nil {
			fmt.Printf("Failed to update rgsgroup: %v\n", err)
		}

	} else if len(rgsprivate) == 0 && msg.MessageType == "private" {
		newPrivateMsgHandle := DSReqModel.RgsPrivate{
			QQID:      int(msg.Sender.UserID),
			SessionID: sessionID,
		}
		//fmt.Println(responseBody.Output.SessionID)
		err = DBControlApi.Db.Create(&newPrivateMsgHandle, "localrgsprivate")
		if err != nil {
			fmt.Printf("Failed to update rgsprivate: %v\n", err)
		}

	}

	if handler, exists := HTTPReq.ReqApiMap[SendMsgtype]; exists {
		if SendMsgtype == ReqApiConst.SEND_GROUP_MSG {
			handler(SendMsgtype, MessageModel.NormalRespMessage(sendMsgID, Resp))
		} else if SendMsgtype == ReqApiConst.SEND_PRIVATE_MSG {
			handler(SendMsgtype, MessageModel.NormalPrivateRespMessage(sendMsgID, Resp))
		}

	}
	return

}
func (n *DeepSeekManageHandle) HandleCloudDeepseekMessage(msg MessageModel.Message) {

	var rgsgroup []DSReqModel.RgsGroup
	var rgsprivate []DSReqModel.RgsPrivate
	var sessionID string
	var SendMsgtype string
	var sendMsgID int64
	var RobReqMsg string

	if msg.MessageType == "private" {
		_, err := DBControlApi.Db.Where("rgsprivate", &rgsprivate, "QQID = ?", msg.Sender.UserID)
		if err != nil {
			sessionID = ""
		} else {
			if len(rgsgroup) == 0 {
				sessionID = ""
			} else {
				sessionID = rgsgroup[0].SeessionID
			}

		}
		//设置发送的群或者人
		sendMsgID = msg.Sender.UserID
		//设置该调用的类型
		SendMsgtype = ReqApiConst.SEND_PRIVATE_MSG
		//定制送个rob消息类型
		RobReqMsg = msg.RawMessage
	} else if msg.MessageType == "group" {
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

		sendMsgID = msg.GroupID
		SendMsgtype = ReqApiConst.SEND_GROUP_MSG
		RobReqMsg = Tool.ExtractMessageForRob(msg.RawMessage)
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
			"prompt":     Tool.Int64toString(msg.Sender.UserID) + "-" + msg.Sender.NickName + "：" + RobReqMsg,
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

	if len(rgsgroup) == 0 && msg.MessageType == "group" {
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

	} else if len(rgsprivate) == 0 && msg.MessageType == "private" {
		newPrivateMsgHandle := DSReqModel.RgsPrivate{
			QQID:      int(msg.Sender.UserID),
			SessionID: responseBody.Output.SessionID,
		}
		fmt.Println(responseBody.Output.SessionID)
		err = DBControlApi.Db.Create(&newPrivateMsgHandle, "rgsprivate")
		if err != nil {
			fmt.Printf("Failed to update rgsprivate: %v\n", err)
		}

	}
	if handler, exists := HTTPReq.ReqApiMap[SendMsgtype]; exists {
		if SendMsgtype == ReqApiConst.SEND_GROUP_MSG {
			handler(SendMsgtype, MessageModel.NormalRespMessage(sendMsgID, responseBody.Output.Text))
		} else if SendMsgtype == ReqApiConst.SEND_PRIVATE_MSG {
			handler(SendMsgtype, MessageModel.NormalPrivateRespMessage(sendMsgID, responseBody.Output.Text))
		}

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

// GenerateSessionID 生成一个随机的 sessionId（长度 16 字节 = 32 个十六进制字符）
func GenerateSessionID() string {
	bytes := make([]byte, 16) // 16 字节 = 128 位
	_, err := rand.Read(bytes)
	if err != nil {
		panic("无法生成随机 sessionId")
	}
	return hex.EncodeToString(bytes) // 转换为十六进制字符串
}
