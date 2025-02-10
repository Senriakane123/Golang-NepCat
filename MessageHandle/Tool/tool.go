package Tool

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func AtQQNumber(str string) (bool, []string) {
	//// 正则表达式匹配 [CQ:at,qq=3666859102] 里的 QQ 号
	//re := regexp.MustCompile(`\[CQ:at,qq=(\d+)\]`)
	//
	//// 查找匹配的 QQ 号
	//match := re.FindStringSubmatch(str)
	//if len(match) > 1 {
	//	fmt.Println("提取的 QQ 号:", match[1])
	//	return true, match
	//} else {
	//	fmt.Println("未找到匹配的 QQ 号")
	//	return false, nil
	//}

	// 正则表达式匹配所有 [CQ:at,qq=xxx] 里的 QQ 号
	re := regexp.MustCompile(`\[CQ:at,qq=(\d+)\]`)

	// 查找所有匹配的 QQ 号
	matches := re.FindAllStringSubmatch(str, -1)

	// 存储提取出的 QQ 号
	var QQNumbers []string

	for _, match := range matches {
		if len(match) > 1 {
			QQNumbers = append(QQNumbers, match[1]) // 只存 QQ 号
		}
	}

	if len(QQNumbers) > 0 {
		return true, QQNumbers
	}

	return false, nil
}

func Int64toString(ints int64) string {
	return strconv.FormatInt(ints, 10)
	//msg := "[CQ:at,qq=" + strconv.FormatInt(message.Sender.UserID, 10) + "] 请求格式错误或用户不具备管理员权限"
}

func StringToInt64(str string) int64 {
	// 将字符串转换为 int64
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		return num
	}
	return num
}

// 生成带 CQ:at 和列表的消息
func BuildReplyMessage(Message []string) string {
	var builder strings.Builder

	//// 添加 @ 用户
	//builder.WriteString("[CQ:at,qq=")
	//builder.WriteString(strconv.FormatInt(userID, 10))
	//builder.WriteString("]\n")

	// 添加列表内容
	// 遍历 Message 列表并动态添加
	for _, item := range Message {
		builder.WriteString(item)
		builder.WriteString("\n")
	}

	return builder.String()
}

// 服务类型检查
func ParseServiceCommand(message string) []string {
	// 正则匹配服务编号，例如 "服务1" 或 "服务5.1"
	re := regexp.MustCompile(`服务(\d+(\.\d+)?)`)
	matches := re.FindStringSubmatch(message)

	// 如果匹配到服务编号并且有小数点，去掉小数点部分
	if len(matches) > 1 {
		// 去掉小数点后面的部分
		if strings.Contains(matches[1], ".") {
			matches[1] = strings.Split(matches[1], ".")[0]
		}

		matches[2] = strings.ReplaceAll(matches[2], ".", "")
	}

	return matches
}

// CheckBanFormat 检查字符串是否符合禁言格式
func CheckBanFormat(input string) (bool, []int64, int) {
	// 正则匹配所有 "[CQ:at,qq=xxx]" 形式的QQ号
	qqPattern := regexp.MustCompile(`\[CQ:at,qq=(\d+)\]`)
	qqMatches := qqPattern.FindAllStringSubmatch(input, -1)

	// 正则匹配禁言时长（形如 `-60`）
	timePattern := regexp.MustCompile(`-(\d+)$`)
	timeMatch := timePattern.FindStringSubmatch(input)

	// 如果没有匹配到QQ号或者没有禁言时长，返回 false
	if len(qqMatches) == 0 || len(timeMatch) == 0 {
		return false, nil, 0
	}

	// 解析所有的QQ号
	var qqNumbers []int64
	for _, match := range qqMatches {
		if len(match) > 1 {
			qq, err := strconv.ParseInt(match[1], 10, 64)
			if err == nil {
				qqNumbers = append(qqNumbers, qq)
			}
		}
	}

	// 解析禁言时长
	banTime, err := strconv.Atoi(timeMatch[1])
	if err != nil {
		return false, nil, 0
	}

	return true, qqNumbers, banTime
}

// CheckBanFormat 检查字符串是否符合禁言格式
func CheckKickFormat(input string) (bool, []int64, int) {
	// 正则匹配所有 "[CQ:at,qq=xxx]" 形式的QQ号
	qqPattern := regexp.MustCompile(`\[CQ:at,qq=(\d+)\]`)
	qqMatches := qqPattern.FindAllStringSubmatch(input, -1)

	// 正则匹配禁言时长（形如 `-60`）
	timePattern := regexp.MustCompile(`-(\d+)$`)
	timeMatch := timePattern.FindStringSubmatch(input)

	// 如果没有匹配到QQ号或者没有禁言时长，返回 false
	if len(qqMatches) == 0 || len(timeMatch) == 0 {
		return false, nil, 0
	}

	// 解析所有的QQ号
	var qqNumbers []int64
	for _, match := range qqMatches {
		if len(match) > 1 {
			qq, err := strconv.ParseInt(match[1], 10, 64)
			if err == nil {
				qqNumbers = append(qqNumbers, qq)
			}
		}
	}

	// 解析禁言时长
	banTime, err := strconv.Atoi(timeMatch[1])
	if err != nil {
		return false, nil, 0
	}

	return true, qqNumbers, banTime
}

// ExtractTags 解析输入格式，提取图片数量和 tag 关键词
func ExtractTags(input string) (int, []string, error) {
	// 正则匹配 `Tag涩图-数字-tag1,tag2`
	re := regexp.MustCompile(`Tag涩图-(\d+)-(.+)`)

	// 查找匹配项
	matches := re.FindStringSubmatch(input)
	if matches == nil || len(matches) < 3 {
		return 0, nil, fmt.Errorf("格式不正确")
	}

	// 解析图片数量
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, nil, fmt.Errorf("图片数量解析失败")
	}

	// 解析 tag，按中文 `，` 或英文 `,` 分割
	tags := strings.FieldsFunc(matches[2], func(r rune) bool {
		return r == '，' || r == ',' // 兼容中英文逗号
	})

	return num, tags, nil
}

// 将图片文件转换为Base64编码，并根据图片类型添加前缀
func ImageToBase64(filePath string) (string, error) {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容并获取文件类型
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 检测文件类型
	contentType := http.DetectContentType(fileBytes)
	var base64Prefix string

	// 根据文件类型选择合适的前缀
	switch contentType {
	case "image/jpeg":
		base64Prefix = "data:image/jpeg;base64,"
	case "image/png":
		base64Prefix = "data:image/png;base64,"
	case "image/gif":
		base64Prefix = "data:image/gif;base64,"
	default:
		return "", fmt.Errorf("不支持的图片格式: %s", contentType)
	}

	// 将文件内容编码为Base64
	base64Str := base64.StdEncoding.EncodeToString(fileBytes)

	// 添加前缀并返回
	return base64Prefix + base64Str, nil
}

func ChangePicMD(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString("\n") // 添加不可见字符
	fmt.Println("图片数据已修改")
	return err
}

func ExtractPetIdNumber(text string) (string, error) {
	// 编译正则表达式，匹配 "用户注册-" 后面的数字（包括负数）
	re := regexp.MustCompile(`用户注册-?(\d+)`)
	matches := re.FindStringSubmatch(text)

	// 检查是否匹配成功
	if len(matches) < 2 {
		return "", fmt.Errorf("未找到匹配的数字")
	}

	// 返回匹配的数字
	return matches[1], nil
}
