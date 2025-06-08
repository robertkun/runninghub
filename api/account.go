package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AccountStatusResponse 账户状态响应
type AccountStatusResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		RemainCoins      string `json:"remainCoins"`
		CurrentTaskCounts string `json:"currentTaskCounts"`
	} `json:"data"`
}

// GetAccountStatus 获取账户信息
func GetAccountStatus(apiKey string) (*AccountStatusResponse, error) {
	url := "https://www.runninghub.cn/uc/openapi/accountStatus"
	
	// 构建请求体
	reqBody := map[string]string{
		"apikey": apiKey,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("构建请求体失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "www.runninghub.cn")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result AccountStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
} 