package ResHandle

import (
	"NepcatGoApiReq/MessageHandle"
	"NepcatGoApiReq/MessageModel"
	"encoding/json"
	"fmt"
	"regexp"
)

// 解析 JSON 并处理
func HandleMessage(jsonData string) {
	var msg MessageModel.Message
	err := json.Unmarshal([]byte(jsonData), &msg)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}

	switch msg.PostType {
	case "meta_event":
		handleMetaEvent(msg)
	case "message":
		handleChatMessage(msg)
	default:
		fmt.Println("未知消息类型:", msg.PostType)
	}
}

// 处理 meta_event，例如心跳包
func handleMetaEvent(msg MessageModel.Message) {
	if msg.MetaEventType == "heartbeat" {
		fmt.Println("收到心跳包")
	} else if msg.MetaEventType == "lifecycle" {
		fmt.Println("机器人上线")
	}
}

// 处理聊天消息
func handleChatMessage(msg MessageModel.Message) {
	fmt.Printf("收到来自用户 %d 的消息: %s\n", msg.UserID, msg.RawMessage)
	switch msg.MessageType {
	case "group":
		MessageHandle.GroupMessageHandle(msg)
		fmt.Println("暂定调用http回复")
	case "private":
		fmt.Println("消息，暂无处理")
	}
	//if msg.MessageType == "group" {
	//	fmt.Printf("群聊消息，群ID: %d\n", msg.GroupID)
	//	// 可进一步处理 at 消息、图片等
	//}
}

func AtQQNumber(str string) (bool, []string) {
	// 正则表达式匹配 [CQ:at,qq=3666859102] 里的 QQ 号
	re := regexp.MustCompile(`\[CQ:at,qq=(\d+)\]`)

	// 查找匹配的 QQ 号
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		fmt.Println("提取的 QQ 号:", match[1])
		return true, match
	} else {
		fmt.Println("未找到匹配的 QQ 号")
		return false, nil
	}
}
