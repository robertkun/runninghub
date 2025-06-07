package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type UploadResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		FileName string `json:"fileName"`
		FileType string `json:"fileType"`
	} `json:"data"`
}

// UploadImage 上传图片到 RunningHub 服务器
// filePath: 本地图片文件路径
func UploadImage(filePath string) (*UploadResponse, error) {
	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加apiKey字段
	if err := writer.WriteField("apiKey", ApiKey); err != nil {
		return nil, fmt.Errorf("写入apiKey失败: %v", err)
	}

	// 添加fileType字段
	if err := writer.WriteField("fileType", "image"); err != nil {
		return nil, fmt.Errorf("写入fileType失败: %v", err)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 创建文件表单字段
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建文件表单字段失败: %v", err)
	}

	// 复制文件内容到表单
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 关闭writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭writer失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", "https://www.runninghub.cn/task/openapi/upload", body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Host", "www.runninghub.cn")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var uploadResp UploadResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &uploadResp, nil
}
