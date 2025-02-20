package MemoryIDCtl

import (
	"NepcatGoApiReq/MessageHandle/DeepSeekReqHandle/DSReqModel"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CreateMemoryResponse 定义创建 API 响应结构
//type CreateMemoryResponse struct {
//	MemoryID string `json:"memoryId"`
//}

// ListMemoryResponse 定义获取长期记忆体列表的响应结构
type ListMemoryResponse struct {
	Memories []struct {
		Description string `json:"description"`
		MemoryID    string `json:"memoryId"`
	} `json:"memories"`
	NextToken string `json:"nextToken,omitempty"`
}

// MemoryIDHandle 结构体
type MemoryIDHandle struct {
	MemoryID string
	APIKey   string
}

// 创建长期记忆体1
func (n *MemoryIDHandle) CreateMemory(Requrl, workspaceId, apiKey string) error {
	// 构造请求 URL
	url := Requrl + fmt.Sprintf("/%s/memories", workspaceId)

	// 创建空的请求体
	reqBody := []byte("{}")

	// 创建新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("请求创建失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	resp, err := http.DefaultClient.Do(req)
	fmt.Println(resp)
	if err != nil {
		return fmt.Errorf("请求错误: %v", err)
	}
	defer resp.Body.Close()

	// 打印响应状态码
	fmt.Println("Response Status:", resp.Status)

	// 如果状态码不是 200，打印详细的错误信息
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("创建 Memory 失败, 状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	// 解析响应内容
	var result DSReqModel.CreateMemoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析 Memory 响应失败: %v", err)
	}

	// 设置 MemoryID
	n.MemoryID = result.MemoryID
	return nil
}

// DeleteMemory 删除长期记忆体
func (n *MemoryIDHandle) DeleteMemory(workspaceID, apiKey, memoryID string) error {
	url := fmt.Sprintf("https://api.your-service.com/%s/memories/%s", workspaceID, memoryID) // 替换 API 地址

	// 创建 HTTP 请求
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("创建删除请求失败: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求错误: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("删除 Memory 失败: %s", string(body))
	}

	return nil
}

// ListMemories 获取长期记忆体列表
func (n *MemoryIDHandle) ListMemories(workspaceID string, maxResults int, nextToken string, apiKey string) (*ListMemoryResponse, error) {
	url := fmt.Sprintf("https://api.your-service.com/%s/memories?maxResults=%d", workspaceID, maxResults) // 替换 API 地址
	if nextToken != "" {
		url += "&nextToken=" + nextToken
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求错误: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取 Memory 列表失败: %s", string(body))
	}

	// 解析响应
	var result ListMemoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 Memory 列表响应失败: %v", err)
	}

	return &result, nil
}
