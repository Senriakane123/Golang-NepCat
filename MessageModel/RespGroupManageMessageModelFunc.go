package MessageModel

import (
	"fmt"
	"strings"
)

func NormalRespMessage(GroupId int64, RespMessage string) map[string]interface{} {
	// 构造返回数据

	// 构造返回数据
	message := map[string]interface{}{
		"group_id": GroupId, // 群号
		//"message": "[CQ:at,qq="+Int64toString(UserId)+"]",
		"message": RespMessage,
	}

	return message
}

func NormalRespMessageType2(GroupId int64, RespMessage []string) map[string]interface{} {

	var builder strings.Builder
	for _, item := range RespMessage {
		builder.WriteString(item)
		builder.WriteString("\n")
	}

	//return builder.String()

	messages := []map[string]interface{}{}

	messages = append(messages, map[string]interface{}{
		"type": "text",
		"data": map[string]interface{}{
			"text": builder.String(),
		},
	})

	message := map[string]interface{}{
		"group_id": GroupId,  // 替换为你的QQ群号
		"message":  messages, // 发送的消息内容
	}

	return message
}

func BanRespMessage(GroupId int64, UserId int64, time int64) map[string]interface{} {
	message := map[string]interface{}{
		"group_id": GroupId, // 替换为你的QQ群号
		"user_id":  UserId,  // 发送的消息内容

		"duration": time, // 发送的消息内容
	}
	return message
}
func KickRespMessage(GroupId int64, UserId int64, kickboolen bool) map[string]interface{} {
	message := map[string]interface{}{
		"group_id": GroupId, // 替换为你的QQ群号
		"user_id":  UserId,  // 发送的消息内容

		"reject_add_request": kickboolen, // 发送的消息内容
	}
	return message
}

func GroupBanRespMessage(GroupId int64, BanBoolen bool) map[string]interface{} {
	message := map[string]interface{}{
		"group_id": GroupId,   // 替换为你的QQ群号
		"enable":   BanBoolen, // 发送的消息内容
	}
	return message
}

func SendRandomPic(ImageBase64strs []string, GroupId int64, info *[]PixivImage) map[string]interface{} {
	fmt.Println(info)

	// 消息列表
	messages := []map[string]interface{}{}
	num := 0
	ErrImageNum := 0
	// 遍历所有图片
	for _, img := range *info {
		// 追加文本消息
		if ImageBase64strs[num] == "" {
			ErrImageNum++
			continue
		}
		messages = append(messages, map[string]interface{}{
			"type": "text",
			"data": map[string]interface{}{
				"text": fmt.Sprintf("这是一张来自 Pixiv 的图片，作者: %s，标题: %s, UID：%d，URL：%s",
					img.Author, img.Title, img.PID, img.URLs.Original),
			},
		})
		messages = append(messages, map[string]interface{}{
			"type": "image",
			"data": map[string]interface{}{
				"file": ImageBase64strs[num],
			},
		})
		num++
	}

	messages = append(messages, map[string]interface{}{
		"type": "text",
		"data": map[string]interface{}{
			"text": fmt.Sprintf("获取图片成功：%d，失败：%d",
				len(*info)-ErrImageNum, ErrImageNum),
		},
	})

	// 构造返回数据
	message := map[string]interface{}{
		"group_id": GroupId, // 群号
		//"message": "[CQ:at,qq="+Int64toString(UserId)+"]",
		"message": messages,
	}

	//fmt.Println(message)
	return message
}
