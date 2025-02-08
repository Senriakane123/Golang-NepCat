package MessageModel

var serverMenu = []string{
	"请选择你的服务",
	"请求格式为'服务'+对应的服务编号",
	"例如'服务1'",
	"1 你是一只猫娘",
	"2 你是我的※奴",
	"3 群主是我的rbq",
	"4 尝试接入deepseek（正在开发中）",
	"5 群管理",
	"6 随机涩图 ",
}

var childServerMenu5 = []string{
	"tips：此服务需要拥有管理员权限",
	"以秒为计算，如果填入60则禁言1分钟",
	"1 禁言 禁言格式为 '@Bot禁言-@群友-禁言时间",
	"2 群体禁言 禁言格式为 '@Bot禁言-@群友1@群友2@群友3'-禁言时间",
	"3 全体禁言 格式为 ‘@Bot全体禁言’",
	"4 解除全体禁言 格式为 ‘@Bot解除全体禁言’",
	"5 踢人 格式为'@Bot踢人-@群友'",
}

var childServerMenu6 = []string{
	"一次请求最多请求三张图片",
	"1 随机涩图 （请求格式为'@Bot随机涩图'）",
	"2 Tag涩图 （请求格式为'@BotTag涩图-图片数量-tag'，tag之间用逗号间隔，例如'@Bot随机涩图-2-碧蓝档案，足控'）",
	//"3 图片识别 （请求格式为'@'）"
}

func GetServerList() []string {
	return serverMenu
}

func GetChildServerList5() []string {
	return childServerMenu5
}

func GetChildServerList6() []string {
	return childServerMenu6
}
