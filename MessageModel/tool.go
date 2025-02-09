package MessageModel

import "strconv"

func Int64toString(ints int64) string {
	return strconv.FormatInt(ints, 10)
	//msg := "[CQ:at,qq=" + strconv.FormatInt(message.Sender.UserID, 10) + "] 请求格式错误或用户不具备管理员权限"
}
