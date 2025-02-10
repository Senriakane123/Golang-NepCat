package R18PicManage

import (
	"NepcatGoApiReq/HTTPReq"
	"NepcatGoApiReq/MessageHandle/Tool"
	"NepcatGoApiReq/MessageModel"
	"NepcatGoApiReq/ReqApiConst"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ReqParam struct {
	Num     int      `json:"num,omitempty"`
	R18     int      `json:"r18,omitempty"`
	Keyword string   `json:"keyword,omitempty"`
	Tags    []string `json:"tag,omitempty"`
	Size    []string `json:"size,omitempty"`
	Proxy   string   `json:"proxy,omitempty"`
}

type PicManage struct {
	folderPath string
	maxImages  int
	handler    map[string]func(MessageModel.Message)
	Params     ReqParam
}

// 文件夹路径
//const folderPath = "PhotoGallery"

// 最大允许的图片数量
const maxImages = 100

// 确保文件夹存在
func (obj *PicManage) ensureFolderExists() error {
	if _, err := os.Stat(obj.folderPath); os.IsNotExist(err) {
		return os.Mkdir(obj.folderPath, os.ModePerm)
	}
	return nil
}
func (n *PicManage) HandlerInit() {
	n.maxImages = 10

	n.folderPath = "PhotoGallery" // 默认路径

	var picManagekeywordHandlers = map[string]func(MessageModel.Message){
		"随机涩图":  n.RandomPic,
		"Tag涩图": n.RandomPicByTagOrNum,
		//"随机R18涩图"
	}
	n.handler = picManagekeywordHandlers
}

// **获取按长度排序的关键词**
func (n *PicManage) getSortedKeywords() []string {
	keys := make([]string, 0, len(n.handler))
	for key := range n.handler {
		keys = append(keys, key)
	}
	// 按字符串长度从长到短排序，保证 "解除全体禁言" 先匹配，而不是 "禁言"
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}

//// 关键词映射到处理函数
//var picManagekeywordHandlers = map[string]func(MessageModel.Message){
//	"随机涩图": handleBan,
//}

// 统一处理消息
func (n *PicManage) HandlePicManageMessage(message MessageModel.Message) bool {
	sortedKeywords := n.getSortedKeywords() // 获取按长度排序的关键词
	for _, keyword := range sortedKeywords {
		if strings.HasPrefix(message.RawMessage, keyword) || strings.Contains(message.RawMessage, keyword) {
			handler := n.handler[keyword]
			handler(message)
			return true // 处理完一个就返回，避免重复触发
		}
	}
	return false
}

// 获取图片 URL
func (n *PicManage) fetchImageURL(reqParams *ReqParam) (*[]MessageModel.PixivImage, error) {
	apiURL := "https://api.lolicon.app/setu/v2"

	// 处理默认参数
	if reqParams.Num <= 0 {
		reqParams.Num = 1
	}
	if reqParams.R18 < 0 || reqParams.R18 > 2 {
		reqParams.R18 = 0
	}
	if reqParams.Size == nil {
		reqParams.Size = []string{"original"}
	}

	// 将参数编码为 JSON
	jsonData, err := json.Marshal(reqParams)
	if err != nil {
		return nil, fmt.Errorf("JSON 序列化失败: %v", err)
	}

	// 发送 POST 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 API 失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应数据失败: %v", err)
	}

	// 解析 JSON
	var apiResp MessageModel.APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	// 检查错误信息
	if apiResp.Error != "" {
		return nil, fmt.Errorf("API 返回错误: %s", apiResp.Error)
	}

	// 确保返回的数据不为空
	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("API 返回空数据")
	}

	// 返回第一张图片
	return &apiResp.Data, nil
}

// 下载并保存图片，返回图片base64编码和文件路径
func (n *PicManage) downloadImage(imageURL string, uid int) (error, string) {
	// 生成唯一文件名（时间戳）
	err := n.ensureFolderExists()
	if err != nil {
		fmt.Println("操作文件失败")
		return err, ""
	}
	err = n.cleanFolderIfNeeded()
	if err != nil {
		fmt.Println("清空文件夹失败")
		return err, ""
	}

	fileName := fmt.Sprintf("%d_%s.jpg", uid, time.Now().Format("20060102-150405"))
	filePath := filepath.Join(n.folderPath, fileName)

	// 发送 HTTP 请求获取图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("下载图片失败: %v", err), ""
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode), ""
	}

	// 创建文件并写入数据
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err), ""
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("保存图片失败: %v", err), ""
	}

	fmt.Println("✅ 图片已保存:", filePath)
	err = Tool.ChangePicMD(filePath)
	if err != nil {
		return fmt.Errorf("添加不可见数据失败: %v", err), ""
	}
	// 获取图片的 Base64 编码
	base64Str, err := Tool.ImageToBase64(filePath)
	if err != nil {
		return fmt.Errorf("转换图片为Base64失败: %v", err), ""
	}

	return nil, base64Str
}

// 清理文件夹
func (n *PicManage) cleanFolderIfNeeded() error {
	files, err := os.ReadDir(n.folderPath)
	if err != nil {
		return fmt.Errorf("读取文件夹失败: %v", err)
	}

	if len(files) > maxImages {
		fmt.Println("⚠️ 超过最大图片数，清空文件夹...")
		for _, file := range files {
			err := os.Remove(filepath.Join(n.folderPath, file.Name()))
			if err != nil {
				fmt.Println("删除失败:", err)
			}
		}
		fmt.Println("✅ 文件夹已清空")
	}
	return nil
}

// 随机涩图
func (n *PicManage) RandomPic(message MessageModel.Message) {
	var Params ReqParam
	//var path string
	//var Paths []string
	var ImageBase64Str string
	var ImageBase64Strs []string
	Params.R18 = 0
	Params.Tags = nil
	Params.Num = 1
	PicInfo, err := n.fetchImageURL(&Params)
	if err != nil {
		return
	}
	for _, s := range *PicInfo {
		err, ImageBase64Str = n.downloadImage(s.URLs.Original, s.PID)
		//Paths = append(Paths, path)
		ImageBase64Strs = append(ImageBase64Strs, ImageBase64Str)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, "禁言请求", message.GroupID, MessageModel.SendRandomPic(ImageBase64Strs, message.GroupID, PicInfo))
	}

}

// 随机涩图
func (n *PicManage) RandomPicByTagOrNum(message MessageModel.Message) {

	_, Tags, err := Tool.ExtractTags(message.RawMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	var Params ReqParam
	var ImageBase64Str string
	var ImageBase64Strs []string
	Params.R18 = 0
	Params.Tags = Tags
	Params.Num = 1

	PicInfo, err := n.fetchImageURL(&Params)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, s := range *PicInfo {
		err, ImageBase64Str = n.downloadImage(s.URLs.Original, s.PID)
		//Paths = append(Paths, path)
		ImageBase64Strs = append(ImageBase64Strs, ImageBase64Str)
	}
	if handler, exists := HTTPReq.ReqApiMap[ReqApiConst.SEND_GROUP_MSG]; exists {
		handler(ReqApiConst.SEND_GROUP_MSG, "禁言请求", message.GroupID, MessageModel.SendRandomPic(ImageBase64Strs, message.GroupID, PicInfo))
	}
}

func (n *PicManage) PicReqRawMessageConfig(message MessageModel.Message) {

}
