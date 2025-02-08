package MessageModel

import (
	"fmt"
)

func NormalRespMessage(GroupId int64, RespMessage string) map[string]interface{} {
	message := map[string]interface{}{
		"group_id": GroupId,     // 替换为你的QQ群号
		"message":  RespMessage, // 发送的消息内容
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

func SendRandomPic(GroupId int64, info *[]PixivImage) map[string]interface{} {
	fmt.Println(info)

	// 消息列表
	messages := []map[string]interface{}{}

	// 遍历所有图片
	for _, img := range *info {
		// 追加文本消息
		messages = append(messages, map[string]interface{}{
			"type": "text",
			"data": map[string]interface{}{
				"text": fmt.Sprintf("这是一张来自 Pixiv 的图片，作者: %s，标题: %s, UID：%d，URL：%s",
					img.Author, img.Title, img.PID, img.URLs.Original),
			},
		})

		// 追加图片消息
		messages = append(messages, map[string]interface{}{
			"type": "image",
			"data": map[string]interface{}{
				"file": img.URLs.Original,
			},
		})
	}

	// 构造返回数据
	message := map[string]interface{}{
		"group_id": GroupId, // 群号
		"message":  messages,
	}

	fmt.Println(message)
	return message
}
