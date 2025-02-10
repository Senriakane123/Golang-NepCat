package HTTPReq

import (
	"NepcatGoApiReq/ReqApiConst"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 定义 API 处理函数类型
type ReqHandler func(ReqMessage string, Resp map[string]interface{})

func HandleMessage(ReqMessage string, Resp map[string]interface{}) {
	url := "http://127.0.0.1:3000/" + ReqMessage
	//var message map[string]interface{}
	//if Resp != nil {
	//	message = Resp
	//} else {
	//	// OneBot 11 服务器地址
	//
	//	// 发送的消息内容
	//	message = map[string]interface{}{
	//		"group_id": GroupId,     // 替换为你的QQ群号
	//		"message":  RespMessage, // 发送的消息内容
	//	}
	//
	//}

	// 将 Go 的 map 转换为 JSON
	jsonData, err := json.Marshal(Resp)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your_access_token") // 如果需要身份验证

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 输出 API 返回结果
	fmt.Println("状态码:", resp.Status)
	fmt.Println("响应数据:", string(body))

	//WebSocket 服务器地址（注意修改为你的 go-cqhttp 监听地址）
}

// 维护一个 API 处理映射
var ReqApiMap = map[string]ReqHandler{}

// HttpReqInit 初始化 API 处理映射
func HttpReqInit(ReqType string, handler ReqHandler) {
	if handler != nil {
		handler = handler
	} else {
		handler = nil
	}
	ReqApiMap[ReqType] = handler
}

// 统一初始化所有 API
func InitAllApis() {
	ReqTypes := []string{
		// OneBot 11 API
		ReqApiConst.SEND_PRIVATE_MSG,
		ReqApiConst.SEND_GROUP_MSG,
		ReqApiConst.SEND_MSG,
		ReqApiConst.DELETE_MSG,
		ReqApiConst.GET_MSG,
		ReqApiConst.GET_FORWARD_MSG,
		ReqApiConst.SEND_LIKE,
		ReqApiConst.SET_GROUP_KICK,
		ReqApiConst.SET_GROUP_BAN,
		ReqApiConst.SET_GROUP_WHOLE_BAN,
		ReqApiConst.SET_GROUP_ADMIN,
		ReqApiConst.SET_GROUP_CARD,
		ReqApiConst.SET_GROUP_NAME,
		ReqApiConst.SET_GROUP_LEAVE,
		ReqApiConst.SET_GROUP_SPECIAL_TITLE,
		ReqApiConst.SET_FRIEND_ADD_REQUEST,
		ReqApiConst.SET_GROUP_ADD_REQUEST,
		ReqApiConst.GET_LOGIN_INFO,
		ReqApiConst.GET_STRANGER_INFO,
		ReqApiConst.GET_FRIEND_LIST,
		ReqApiConst.GET_GROUP_INFO,
		ReqApiConst.GET_GROUP_LIST,
		ReqApiConst.GET_GROUP_MEMBER_INFO,
		ReqApiConst.GET_GROUP_MEMBER_LIST,
		ReqApiConst.GET_GROUP_HONOR_INFO,
		ReqApiConst.GET_COOKIES,
		ReqApiConst.GET_CSRF_TOKEN,
		ReqApiConst.GET_CREDENTIALS,
		ReqApiConst.GET_RECORD,
		ReqApiConst.GET_IMAGE,
		ReqApiConst.CAN_SEND_IMAGE,
		ReqApiConst.CAN_SEND_RECORD,
		ReqApiConst.GET_STATUS,
		ReqApiConst.GET_VERSION_INFO,
		ReqApiConst.CLEAN_CACHE,

		// go-cqhttp API
		ReqApiConst.SET_QQ_PROFILE,
		ReqApiConst.GET_MODEL_SHOW,
		ReqApiConst.SET_MODEL_SHOW,
		ReqApiConst.GET_ONLINE_CLIENTS,
		ReqApiConst.DELETE_FRIEND,
		ReqApiConst.MARK_MSG_AS_READ,
		ReqApiConst.SEND_GROUP_FORWARD_MSG,
		ReqApiConst.SEND_PRIVATE_FORWARD_MSG,
		ReqApiConst.GET_GROUP_MSG_HISTORY,
		ReqApiConst.OCR_IMAGE,
		ReqApiConst.GET_GROUP_SYSTEM_MSG,
		ReqApiConst.GET_ESSENCE_MSG_LIST,
		ReqApiConst.GET_GROUP_AT_ALL_REMAIN,
		ReqApiConst.SET_GROUP_PORTRAIT,
		ReqApiConst.SET_ESSENCE_MSG,
		ReqApiConst.DELETE_ESSENCE_MSG,
		ReqApiConst.SEND_GROUP_SIGN,
		ReqApiConst.SEND_GROUP_NOTICE,
		ReqApiConst.GET_GROUP_NOTICE,
		ReqApiConst.UPLOAD_GROUP_FILE,
		ReqApiConst.DELETE_GROUP_FILE,
		ReqApiConst.CREATE_GROUP_FILE_FOLDER,
		ReqApiConst.DELETE_GROUP_FOLDER,
		ReqApiConst.GET_GROUP_FILE_SYSTEM_INFO,
		ReqApiConst.GET_GROUP_ROOT_FILES,
		ReqApiConst.GET_GROUP_FILES_BY_FOLDER,
		ReqApiConst.GET_GROUP_FILE_URL,
		ReqApiConst.UPLOAD_PRIVATE_FILE,
		ReqApiConst.DOWNLOAD_FILE,
		ReqApiConst.CHECK_URL_SAFELY,
		ReqApiConst.HANDLE_QUICK_OPERATION,

		// napcat API
		ReqApiConst.SET_GROUP_SIGN,
		ReqApiConst.ARK_SHARE_PEER,
		ReqApiConst.ARK_SHARE_GROUP,
		ReqApiConst.GET_ROBOT_UIN_RANGE,
		ReqApiConst.SET_ONLINE_STATUS,
		ReqApiConst.GET_FRIENDS_WITH_CATEGORY,
		ReqApiConst.SET_QQ_AVATAR,
		ReqApiConst.GET_FILE,
		ReqApiConst.FORWARD_FRIEND_SINGLE_MSG,
		ReqApiConst.FORWARD_GROUP_SINGLE_MSG,
		ReqApiConst.TRANSLATE_EN2ZH,
		ReqApiConst.SET_MSG_EMOJI_LIKE,
		ReqApiConst.SEND_FORWARD_MSG,
		ReqApiConst.MARK_PRIVATE_MSG_AS_READ,
		ReqApiConst.MARK_GROUP_MSG_AS_READ,
		ReqApiConst.GET_FRIEND_MSG_HISTORY,
		ReqApiConst.CREATE_COLLECTION,
		ReqApiConst.GET_COLLECTION_LIST,
		ReqApiConst.SET_SELF_LONGNICK,
		ReqApiConst.GET_RECENT_CONTACT,
		ReqApiConst.MARK_ALL_AS_READ,
		ReqApiConst.GET_PROFILE_LIKE,
		ReqApiConst.FETCH_CUSTOM_FACE,
		ReqApiConst.FETCH_EMOJI_LIKE,
		ReqApiConst.SET_INPUT_STATUS,
		ReqApiConst.GET_GROUP_INFO_EX,
		ReqApiConst.GET_GROUP_IGNORE_ADD_REQUEST,
		ReqApiConst.DEL_GROUP_NOTICE,
		ReqApiConst.FRIEND_POKE,
		ReqApiConst.GROUP_POKE,
		ReqApiConst.NC_GET_PACKET_STATUS,
		ReqApiConst.NC_GET_USER_STATUS,
		ReqApiConst.NC_GET_RKEY,
		ReqApiConst.GET_GROUP_SHUT_LIST,
		ReqApiConst.GET_MINI_APP_ARK,
		ReqApiConst.GET_AI_RECORD,
		ReqApiConst.GET_AI_CHARACTERS,
		ReqApiConst.SEND_GROUP_AI_RECORD,
		ReqApiConst.SEND_POKE,
	}

	for _, api := range ReqTypes {
		HttpReqInit(api, HandleMessage) // 绑定默认处理
	}
}

//func HttpReqHandle(ReqType string) {
//
//	// 调用请求
//	if handler, exists := ReqApiMap["SEND_PRIVATE_MSG"]; exists {
//		handler("Hello QQ!")
//	}
//}
