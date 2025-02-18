package NepCatSysHttpReq

import (
	"NepcatGoApiReq/NepCatSysHttpReq/NepcatReqModel"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NepCatHttpReq struct {
	Baseurl       string        `json:"baseurl"`
	PasswordToken LoginRequest  `json:"token"`
	BearerToken   LoginResponse `json:"bearerToken"`
	Config        ConfigRequest `json:"config"`
}

const baseURL = "http://127.0.0.1:6099"

type ConfigRequest struct {
	Config string `json:"config"`
}

// 定义登录请求结构体
type LoginRequest struct {
	Token string `json:"token"`
}

// 定义登录响应结构体（如果需要解析返回的 Token）
type LoginResponse struct {
	BearerToken string `json:"token"`
}

func (n *NepCatHttpReq) BaseUrlAndPasswordinit(Password, Baseurl string) {
	n.Baseurl = Baseurl
	n.PasswordToken.Token = Password
}

func (n *NepCatHttpReq) Login() error {
	// 1. 构造请求体
	loginData := n.PasswordToken
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("json 编码错误: %v", err)
	}

	// 2. 发送 HTTP POST 请求
	url := baseURL + "/api/auth/MusicGetLogic"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 3. 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 5. 解析 JSON 响应
	var loginResp NepcatReqModel.AuthResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return fmt.Errorf("解析 JSON 失败: %v", err)
	}

	fmt.Println("登录成功，获取到 Token:", loginResp.Data.Credential)
	n.BearerToken.BearerToken = loginResp.Data.Credential
	return nil
}

// 定义配置请求结构体

func (n *NepCatHttpReq) SetConfig() error {
	// 1. 构造请求体
	configData := n.Config

	jsonData, err := json.Marshal(configData)
	if err != nil {
		return fmt.Errorf("json 编码错误: %v", err)
	}

	// 2. 发送 HTTP POST 请求
	url := baseURL + "/api/OB11Config/SetConfig"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 3. 设置请求头，使用从登录请求获取的 Token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.BearerToken.BearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Println("配置成功:", string(body))
	return nil
}

// 获取配置的函数
func (n *NepCatHttpReq) GetConfig() (string, error) {
	// 1. 构造请求 URL
	url := baseURL + "/api/OB11Config/GetConfig"

	// 2. 构造请求体（为空）
	jsonData, _ := json.Marshal(map[string]string{})

	// 3. 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 4. 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.BearerToken.BearerToken)

	// 5. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 6. 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 7. 输出返回的 JSON 数据
	fmt.Println("获取到的配置:")
	fmt.Println(string(body))
	return string(body), nil
}
