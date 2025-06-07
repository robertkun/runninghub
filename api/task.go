package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NodeInfo struct {
	NodeId     string      `json:"nodeId"`
	FieldName  string      `json:"fieldName"`
	FieldValue interface{} `json:"fieldValue"`
}

type TaskCreateResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		NetWssUrl  string `json:"netWssUrl"`
		TaskId     string `json:"taskId"`
		ClientId   string `json:"clientId"`
		TaskStatus string `json:"taskStatus"`
		PromptTips string `json:"promptTips"`
	} `json:"data"`
}

type TaskStatusResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"` // 任务状态: QUEUED, RUNNING, FAILED, SUCCESS
}

type TaskOutput struct {
	FileUrl      string `json:"fileUrl"`
	FileType     string `json:"fileType"`
	TaskCostTime string `json:"taskCostTime"`
	NodeId       string `json:"nodeId"`
}

type TaskOutputResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data []TaskOutput `json:"data"`
}

type CancelTaskResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// CreateAdvancedTask 发起高级 ComfyUI 任务
// workflowId: 工作流ID
// nodeInfoList: 节点参数修改列表
func CreateAdvancedTask(workflowId string, nodeInfoList []NodeInfo) (*TaskCreateResponse, error) {
	url := "https://www.runninghub.cn/task/openapi/create"
	method := "POST"

	payload := map[string]interface{}{
		"apiKey":       ApiKey,
		"workflowId":   workflowId,
		"nodeInfoList": nodeInfoList,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 打印请求参数
	fmt.Println("[CreateAdvancedTask] 请求URL:", url)
	fmt.Println("[CreateAdvancedTask] 请求参数:", string(jsonData))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Host", "www.runninghub.cn")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 打印响应内容
	fmt.Println("[CreateAdvancedTask] 响应内容:", string(body))

	var taskResp TaskCreateResponse
	if err := json.Unmarshal(body, &taskResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &taskResp, nil
}

// QueryTaskStatus 查询任务状态
// taskId: 任务ID
func QueryTaskStatus(taskId string) (*TaskStatusResponse, error) {
	url := "https://www.runninghub.cn/task/openapi/status"
	method := "POST"

	payload := map[string]string{
		"apiKey": ApiKey,
		"taskId": taskId,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Host", "www.runninghub.cn")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var statusResp TaskStatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &statusResp, nil
}

// QueryTaskOutputs 查询任务生成结果
// taskId: 任务ID
func QueryTaskOutputs(taskId string) (*TaskOutputResponse, error) {
	url := "https://www.runninghub.cn/task/openapi/outputs"
	method := "POST"

	payload := map[string]string{
		"apiKey": ApiKey,
		"taskId": taskId,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Host", "www.runninghub.cn")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var outputResp TaskOutputResponse
	if err := json.Unmarshal(body, &outputResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &outputResp, nil
}

// CancelTask 取消任务
// taskId: 任务ID
func CancelTask(taskId string) (*CancelTaskResponse, error) {
	url := "https://www.runninghub.cn/task/openapi/cancel"
	method := "POST"

	payload := map[string]string{
		"apiKey": ApiKey,
		"taskId": taskId,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Host", "www.runninghub.cn")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var cancelResp CancelTaskResponse
	if err := json.Unmarshal(body, &cancelResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &cancelResp, nil
}
