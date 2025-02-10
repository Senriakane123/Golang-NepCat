package GameLogic

import (
	"NepcatGoApiReq/Database/DBControlApi"
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/RPGGameHandle/GameDatamodel"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"fmt"
	"log"
	"sort"
	"strings"
)

type GameManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *GameManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"用户注册":     userRegister,
		"获取宠物信息":   getPetInfo,
		"获取注册宠物列表": getEnableRegisPetList,
		"等级查询":     levelQuery,
		"每日签到":     dailySignIn,
	}
	n.Handler = groupManagekeywordHandlers
}

// **获取按长度排序的关键词**
func (n *GameManageHandle) getSortedKeywords() []string {
	keys := make([]string, 0, len(n.Handler))
	for key := range n.Handler {
		keys = append(keys, key)
	}
	// 按字符串长度从长到短排序，保证 "解除全体禁言" 先匹配，而不是 "禁言"
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}

// 统一处理消息
func (n *GameManageHandle) HandleGameManageMessage(message MessageModel.Message) bool {
	sortedKeywords := n.getSortedKeywords() // 获取按长度排序的关键词
	for _, keyword := range sortedKeywords {
		if strings.HasPrefix(message.RawMessage, keyword) || strings.Contains(message.RawMessage, keyword) {
			handler := n.Handler[keyword]
			handler(message)
			return true // 处理完一个就返回，避免重复触发
		}
	}
	return false
}

func userRegister(message MessageModel.Message) {
	fmt.Println("收到用户注册消息:", message)

	// 提取用户 QQ 号
	qqNum := message.Sender.UserID

	// 提取宠物 ID
	PetID, err := Tool.ExtractPetIdNumber(message.RawMessage)
	if err != nil {
		log.Println("提取宠物 ID 失败:", err)
		return
	}

	// 检查用户是否已注册
	var existingUser []GameDatamodel.UserInfo
	_, err = DBControlApi.Db.Where("userinfo", &existingUser, "QQNum = ?", qqNum)
	if err != nil {
		log.Println("查询用户信息失败:", err)
		return
	}

	// 如果用户已存在，直接返回
	if len(existingUser) > 0 {
		fmt.Println("用户已注册，ID:", existingUser[0].ID)
		return
	}

	// 创建新用户
	newUser := GameDatamodel.UserInfo{
		QQNum:   qqNum,
		Name:    message.Sender.NickName,
		Item:    "",
		PetInfo: []GameDatamodel.PersonalPetInfo{},
	}

	// 插入用户数据
	if err = DBControlApi.Db.Create(&newUser, "userinfo"); err != nil {
		log.Println("创建用户失败:", err)
		return
	}

	// 绑定宠物
	newPetInfo := GameDatamodel.PersonalPetInfo{
		UserID:   int64(newUser.ID),
		PetId:    Tool.StringToInt64(PetID),
		Petlevel: 1, // 初始等级
		Exp:      0, // 初始经验
		QQNum:    int(newUser.QQNum),
	}

	// 插入用户宠物数据
	if err = DBControlApi.Db.Create(&newPetInfo, "personalpetinfo"); err != nil {
		log.Println("绑定宠物失败:", err)
		return
	}

	fmt.Println("用户注册成功，ID:", newUser.ID, "绑定宠物 ID:", PetID)

}

func getPetInfo(message MessageModel.Message) {

}
func getEnableRegisPetList(message MessageModel.Message) {
	var EnableRegisPetList []GameDatamodel.Pet
	var EnablePetInfoRegisRespMessage []string
	_, err := DBControlApi.Db.Get(&EnableRegisPetList, "pet")
	if err != nil {
		return
	} else {

		// 遍历宠物列表并格式化

		for _, petInfo := range EnableRegisPetList {
			// 生成宠物信息字符串
			petStr := fmt.Sprintf("ID: %d, \n名称: %s, \n类型: %s, \n基础生命: %d, \n基础攻击: %d, \n基础防御：%d, \n基础能量值：%d \n\n",
				petInfo.ID, petInfo.Name, petInfo.Type, petInfo.BaseHealth, petInfo.BaseAtk, petInfo.BaseDef, petInfo.BaseEnergy)

			EnablePetInfoRegisRespMessage = append(EnablePetInfoRegisRespMessage, petStr)
		}

		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(EnablePetInfoRegisRespMessage)))
		}
	}

}

func levelQuery(message MessageModel.Message) {

}

func dailySignIn(message MessageModel.Message) {

}
