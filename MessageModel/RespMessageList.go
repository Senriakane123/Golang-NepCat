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
	"7 更换机器人头像（需要向机器人所有者获取管理权限，目前是测试开发阶段） ",
	"8 宠物系统（测试开发阶段）",
	"8 DeepSeek（'@Bot接入deepseek'后进入ai聊天模式，再次'@Bot退出deepseek'则会退出AI聊天模式）",
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
	"1 随机涩图 （请求格式为'@Bot随机涩图'）",
	"2 Tag涩图 （请求格式为'@BotTag涩图-图片数量-tag'，tag之间用逗号间隔，例如'@Bot随机涩图-2-碧蓝档案，足控'-----因为服务器响应问题以及保证图片响应速度目前无论选择数目为多少默认为1）",
	//"3 图片识别 （请求格式为'@'）"
}

//var childServerMenu7 = []string{
//	"一次请求最多请求三张图片",
//	"1 随机涩图 （请求格式为'@Bot随机涩图'）",
//	"2 Tag涩图 （请求格式为'@BotTag涩图-图片数量-tag'，tag之间用逗号间隔，例如'@Bot随机涩图-2-碧蓝档案，足控'-----因为服务器响应问题以及保证图片响应速度目前无论选择数目为多少默认为1）",
//	//"3 图片识别 （请求格式为'@'）"
//}

var childServerMenu8 = []string{
	"1 每日签到（@Bot每日签到，第一天100exp，连续签到可额外获得50exp，连续七天达到最大值获得500经验）",
	"2 查询升级需要EXP数值（@bot等级查询-初始等级-预期等级，满级100级）",
	"3 获取宠物信息（@Bot获取宠物信息）",
	"4 注册（注册格式为：@Bot用户注册-宠物ID）",
	"5 BOSS战",
	"6 道具箱",
	"7 获取可注册宠物列表（@Bot获取注册宠物列表）",
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

func GetChildServerList8() []string {
	return childServerMenu8
}
