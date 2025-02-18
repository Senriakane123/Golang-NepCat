package GameLogic

import (
	"NepcatGoApiReq/Database/DBControlApi"
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/RPGGameHandle/GameDatamodel"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

type GameManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *GameManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"用户注册":         n.userRegister,
		"获取宠物信息":     n.getPetInfo,
		"获取注册宠物列表": n.getEnableRegisPetList,
		"等级查询":         n.levelQuery,
		"每日签到":         n.dailySignIn,
		"道具箱":           n.ItemBoxGet,
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

func (n *GameManageHandle) userRegister(message MessageModel.Message) {
	fmt.Println("收到用户注册消息:", message)

	// 提取用户 QQ 号
	qqNum := message.Sender.UserID

	// 提取宠物 ID
	PetID, err := Tool.ExtractPetIdNumber(message.RawMessage)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"用户注册成功，"+"用户请求格式错误"))
		}
		log.Println("提取宠物 ID 失败:", err)
		return
	}

	// 检查用户是否已注册
	var existingUser []GameDatamodel.UserInfo
	_, err = DBControlApi.Db.Where("userinfo", &existingUser, "QQNum = ?", qqNum)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"用户注册成功，"+"查询用户失败"))
		}
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
		QQNum:         qqNum,
		Name:          message.Sender.NickName,
		Item:          "{\"1\":1}",
		PetInfo:       []GameDatamodel.PersonalPetInfo{},
		SignInDayCout: 1,
		SignInTime:    time.Now(), // 赋值当前时间
	}

	// 插入用户数据
	if err = DBControlApi.Db.Create(&newUser, "userinfo"); err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"用户注册成功，"+"创建用户失败"))
		}
		log.Println("创建用户失败:", err)
		return
	}

	// 创建 SkillList
	var SkillList []GameDatamodel.AllSkillList
	SkillList = append(SkillList, GameDatamodel.AllSkillList{
		ID: 1,
	})

	// 将 SkillList 转换为 JSON 字符串
	SkillListJSON, err := json.Marshal(SkillList)
	if err != nil {
		fmt.Println("Error marshaling SkillList:", err)
		return
	}
	// 绑定宠物
	newPetInfo := GameDatamodel.PersonalPetInfo{
		UserID:   int64(newUser.ID),
		PetId:    Tool.StringToInt64(PetID),
		Petlevel: 1, // 初始等级
		Exp:      0, // 初始经验

		QQNum: int(newUser.QQNum),
		Skill: string(SkillListJSON),
	}

	// 插入用户宠物数据
	if err = DBControlApi.Db.Create(&newPetInfo, "personalpetinfo"); err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"用户注册成功，"+"绑定宠物失败"))
		}
		log.Println("绑定宠物失败:", err)
		return
	}
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"用户注册成功，"+fmt.Sprintf("绑定宠物 ID:%d", Tool.StringToInt64(PetID))))
	}

	fmt.Println("用户注册成功，ID:", newUser.ID, "绑定宠物 ID:", PetID)

}

func (n *GameManageHandle) getPetInfo(message MessageModel.Message) {
	var PerPetInfo []GameDatamodel.PersonalPetInfo
	_, err := DBControlApi.Db.Where("personalpetinfo", &PerPetInfo, "QQNum = ?", message.Sender.UserID)
	if err != nil {
		log.Println("查询用户信息失败:", err)
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"查询用户失败请联系管理员"))
		}
		return
	}

	if len(PerPetInfo) == 0 {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"用户暂无注册"))
		}
	}

	///////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////
	// 假设 PerPetInfo[0].Skill 是包含 JSON 字符串的字段
	skillData := []byte(PerPetInfo[0].Skill)

	// 定义 PetSkillList 用来存储解析后的数据
	var PetSkillList []GameDatamodel.AllSkillList

	// 使用 json.Unmarshal 解析字符串
	err = json.Unmarshal(skillData, &PetSkillList)
	if err != nil {
		fmt.Println("Error unmarshalling Skill:", err)
		return
	}

	// 提取 PetSkillList 中的所有 ID
	var skillIDs []int
	for _, skill := range PetSkillList {
		skillIDs = append(skillIDs, skill.ID)
	}

	// 如果没有技能 ID，则不进行查询
	if len(skillIDs) == 0 {
		fmt.Println("No skills found for this pet.")
		return
	}

	// 将 skillIDs 转换为字符串格式，适用于 IN 查询
	// 假设数据库使用的是字符串格式的 ID (你可能需要调整根据数据库的类型)
	var idStrings []string
	for _, id := range skillIDs {
		idStrings = append(idStrings, fmt.Sprintf("%d", id))
	}
	idsStr := strings.Join(idStrings, ",")

	// 查询数据库，使用 IN 子句查询多个技能 ID
	//query := fmt.Sprintf("SELECT * FROM allskilllist WHERE ID IN (%s) AND PetId = ?", idsStr)
	PetSkillList = []GameDatamodel.AllSkillList{}
	// 执行查询
	_, err = DBControlApi.Db.Where("allskilllist", &PetSkillList, "ID = ?", idsStr)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"查询技能失败"))
		}
		fmt.Println("Error querying database:", err)
	} else {
		fmt.Println("Query executed successfully.")
	}

	//PetSkillList = Info.([]GameDatamodel.AllSkillList)

	var perpetskill []string
	for _, skill := range PetSkillList {
		perpetskill = append(perpetskill, skill.SkillName)
	}

	///////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////

	var PetInfo []GameDatamodel.Pet
	_, err = DBControlApi.Db.Where("pet", &PetInfo, "ID = ?", PerPetInfo[0].PetId)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"查询宠物信息失败"))
		}
		log.Println("查询宠物信息失败:", err)
		return
	}

	petStr := fmt.Sprintf("用户的宠物信息为：\n用户QQ:%d\n宠物ID:%d\n宠物名称:%s\n宠物等级:%d\n宠物经验:%d\n宠物技能:%s",
		PerPetInfo[0].QQNum, PetInfo[0].ID, PetInfo[0].Name, PerPetInfo[0].Petlevel, PerPetInfo[0].Exp, strings.Join(perpetskill, ", "))

	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+petStr))
	}
}
func (n *GameManageHandle) getEnableRegisPetList(message MessageModel.Message) {
	var EnableRegisPetList []GameDatamodel.Pet
	var EnablePetInfoRegisRespMessage []string
	_, err := DBControlApi.Db.Get(&EnableRegisPetList, "pet")
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"查询可注册宠物信息失败"))
		}
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

func (n *GameManageHandle) levelQuery(message MessageModel.Message) {

}

func (n *GameManageHandle) dailySignIn(message MessageModel.Message) {
	var Userlist GameDatamodel.UserInfo
	// 提取用户 QQ 号
	qqNum := message.Sender.UserID

	var RespMessage []string
	var ExpItemCode int
	var ExpCardNum int

	itemmap := GameDatamodel.ReturnUserItemList() //item为map[int]int

	// 检查用户是否已注册
	//var existingUser []GameDatamodel.UserInfo
	_, err := DBControlApi.Db.Where("userinfo", &Userlist, "QQNum = ?", qqNum)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"查询用户信息失败"))
		}
		return
	} else {
		// **获取当前时间**
		now := time.Now()
		lastSignIn := Userlist.SignInTime
		fmt.Println(lastSignIn.Year(), now.Year(), lastSignIn.YearDay(), now.YearDay())
		if lastSignIn.YearDay() == now.YearDay() {
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
				handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"用户今日已经签到"))
			}
			return
		}
		// **按天判断是否中断**
		if lastSignIn.Year() != now.Year() || lastSignIn.YearDay() != now.YearDay()-1 {
			Userlist.SignInDayCout = 1 // **断签，重置签到天数**
		} else {
			Userlist.SignInDayCout++ // **连续签到，+1**
			Userlist.SignInTime = time.Now()
		}
		if Userlist.Item != "" {
			err = json.Unmarshal([]byte(Userlist.Item), &itemmap)
			if err != nil {
				return
			}
		}
		switch Userlist.SignInDayCout {
		case 0:
			ExpItemCode = 7
			ExpCardNum = 500
		case 1:
			ExpItemCode = 1
			ExpCardNum = 100
		case 2:
			ExpItemCode = 2
			ExpCardNum = 150
		case 3:
			ExpItemCode = 3
			ExpCardNum = 200
		case 4:
			ExpItemCode = 4
			ExpCardNum = 250
		case 5:
			ExpItemCode = 5
			ExpCardNum = 300
		case 6:
			ExpItemCode = 6
			ExpCardNum = 350
		}

		// **通过 map 快速查找并更新**
		if idx, exists := itemmap[ExpItemCode]; exists {
			itemmap[idx]++
		} else {
			itemmap[ExpItemCode] = 1
		}

		// **转换回 JSON 并存入数据库**
		newItemData, _ := json.Marshal(itemmap)
		Userlist.Item = string(newItemData)

		// **更新数据库**
		if err = DBControlApi.Db.Update(&Userlist, "userinfo"); err != nil {
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
				handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"更新用户信息失败"))
			}
			log.Println("更新用户物品失败:", err)
		}

		RespMessage = append(RespMessage, "签到成功！")
		RespMessage = append(RespMessage, fmt.Sprintf("获得%d经验卡", ExpCardNum))

		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+Tool.BuildReplyMessage(RespMessage)))
		}
	}

}

func (n *GameManageHandle) ItemBoxGet(message MessageModel.Message) {

	var Userlist GameDatamodel.UserInfo
	qqNum := message.Sender.UserID
	itemmap := GameDatamodel.ReturnUserItemList() //item为map[int]int
	var ItemList GameDatamodel.ItemList
	//var UserItemList map[]

	_, err := DBControlApi.Db.Where("userinfo", &Userlist, "QQNum = ?", qqNum)
	if err != nil {
		if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
			handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"查询用户信息失败"))
		}
		return
	} else {

		if Userlist.Item != "" {
			err = json.Unmarshal([]byte(Userlist.Item), &itemmap)
			if err != nil {
				return
			}
		}
		var idStrings []string
		for ID, _ := range itemmap {
			idStrings = append(idStrings, fmt.Sprintf("%d", ID))
		}
		idsStr := strings.Join(idStrings, ",")

		// 执行查询
		_, err = DBControlApi.Db.Where("allskilllist", &ItemList, "ID = ?", idsStr)
		if err != nil {
			if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
				handler(ReqApiConst.SEND_GROUP_MSG, MessageModel.NormalRespMessage(message.GroupID, "[CQ:at,qq="+Tool.Int64toString(message.Sender.UserID)+"]\n"+"查询物品失败"))
			}
			fmt.Println("Error querying database:", err)
		} else {
			fmt.Println("Query executed successfully.")
		}

		//userItems := make(map[int]map[string]int)
		//for _, item := range itemmap {
		//
		//}
	}
}

// 计算等级
func (n *GameManageHandle) LevelCalculate(NowLevel int, exp int, increExp int) (bool, int) {
	baseExp := 100      // 初始基础经验值（1级升2级需要100经验）
	growthFactor := 1.2 // 经验增长系数

	// 计算公式：基础经验 × (增长系数)^当前等级
	requiredExp := float64(baseExp) * math.Pow(growthFactor, float64(NowLevel))
	IntrequiredExp := int(requiredExp)
	if (exp + increExp) >= IntrequiredExp {
		return true, (exp + increExp) - IntrequiredExp
	} else {
		return false, exp + increExp
	}

}

func (n *GameManageHandle) Fight(PetList []GameDatamodel.PersonalPetInfo) {

}
