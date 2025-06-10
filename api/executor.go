package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// WorkflowExecutor 工作流执行器
type WorkflowExecutor struct {
	manager *WorkflowManager
}

// NewWorkflowExecutor 创建工作流执行器
func NewWorkflowExecutor(manager *WorkflowManager) *WorkflowExecutor {
	return &WorkflowExecutor{
		manager: manager,
	}
}

// ExecuteWorkflow 执行工作流
func (we *WorkflowExecutor) ExecuteWorkflow(workflowID string) (*TaskCreateResponse, error) {
	// 获取工作流配置
	config, exists := we.manager.GetWorkflow(workflowID)
	if !exists {
		return nil, fmt.Errorf("工作流不存在: %s", workflowID)
	}

	// 使用工作流配置中的固定参数
	nodeInfoList := make([]NodeInfo, 0, len(config.Params))
	for _, param := range config.Params {
		// 跳过图片输入节点
		if param.IsImage {
			continue
		}
		nodeInfoList = append(nodeInfoList, NodeInfo{
			NodeId:     param.NodeId,
			FieldName:  param.FieldName,
			FieldValue: param.FieldValue,
		})
	}

	// 创建任务
	return CreateAdvancedTask(config.ID, nodeInfoList)
}

// ExecuteWorkflowWithImage 执行带图片的工作流
func (we *WorkflowExecutor) ExecuteWorkflowWithImage(workflowID string, filePath string) (*TaskCreateResponse, error) {
	// 获取工作流配置
	config, exists := we.manager.GetWorkflow(workflowID)
	if !exists {
		return nil, fmt.Errorf("工作流不存在: %s", workflowID)
	}

	// 获取文件类型
	fileType := "image"
	if strings.HasSuffix(strings.ToLower(filePath), ".mp4") {
		fileType = "video"
	}

	// 上传文件
	uploadResp, err := UploadImage(filePath, fileType)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %v", err)
	}

	// 使用工作流配置中的固定参数，但替换图片参数
	nodeInfoList := make([]NodeInfo, 0, len(config.Params))
	for _, param := range config.Params {
		if param.IsImage || param.FieldName == "video" {
			// 设置图片/视频参数
			nodeInfoList = append(nodeInfoList, NodeInfo{
				NodeId:     param.NodeId,
				FieldName:  param.FieldName,
				FieldValue: uploadResp.Data.FileName,
			})
		} else {
			// 设置其他参数
			nodeInfoList = append(nodeInfoList, NodeInfo{
				NodeId:     param.NodeId,
				FieldName:  param.FieldName,
				FieldValue: param.FieldValue,
			})
		}
	}

	// 创建任务
	return CreateAdvancedTask(config.ID, nodeInfoList)
}

// ExecuteWorkflowWithText 执行带文本的工作流
func (we *WorkflowExecutor) ExecuteWorkflowWithText(workflowID string, text string) (*TaskCreateResponse, error) {
	// 获取工作流配置
	config, exists := we.manager.GetWorkflow(workflowID)
	if !exists {
		return nil, fmt.Errorf("工作流不存在: %s", workflowID)
	}

	// 使用工作流配置中的固定参数，但替换文本参数
	nodeInfoList := make([]NodeInfo, 0, len(config.Params))
	for _, param := range config.Params {
		if param.IsImage {
			// 跳过图片输入节点
			continue
		}
		if param.FieldName == "text" {
			// 设置文本参数
			nodeInfoList = append(nodeInfoList, NodeInfo{
				NodeId:     param.NodeId,
				FieldName:  param.FieldName,
				FieldValue: text,
			})
		} else {
			// 设置其他参数
			nodeInfoList = append(nodeInfoList, NodeInfo{
				NodeId:     param.NodeId,
				FieldName:  param.FieldName,
				FieldValue: param.FieldValue,
			})
		}
	}

	// 创建任务
	return CreateAdvancedTask(config.ID, nodeInfoList)
}

// ExecuteWorkflowWithVideoAndAudio 执行带视频和音频的工作流
func (we *WorkflowExecutor) ExecuteWorkflowWithVideoAndAudio(workflowID string, videoPath string, audioPath string) (*TaskCreateResponse, error) {
	// 获取工作流配置
	config, exists := we.manager.GetWorkflow(workflowID)
	if !exists {
		return nil, fmt.Errorf("工作流不存在: %s", workflowID)
	}

	// 上传视频文件
	videoResp, err := UploadImage(videoPath, "image")
	if err != nil {
		return nil, fmt.Errorf("上传视频文件失败: %v", err)
	}
	fmt.Printf("[视频] 上传成功: %s\n", videoResp.Data.FileName)

	// 上传音频文件
	audioResp, err := UploadImage(audioPath, "image")
	if err != nil {
		return nil, fmt.Errorf("上传音频文件失败: %v", err)
	}
	fmt.Printf("[音频] 上传成功: %s\n", audioResp.Data.FileName)

	// 使用工作流配置中的固定参数，但替换视频和音频参数
	nodeInfoList := make([]NodeInfo, 0, len(config.Params))
	for _, param := range config.Params {
		if param.IsImage {
			// 根据节点ID设置不同的文件
			if param.NodeId == "2" {
				// 设置视频参数
				nodeInfoList = append(nodeInfoList, NodeInfo{
					NodeId:     param.NodeId,
					FieldName:  param.FieldName,
					FieldValue: videoResp.Data.FileName,
				})
				fmt.Printf("[视频] 设置到节点 %s: %s\n", param.NodeId, videoResp.Data.FileName)
			} else if param.NodeId == "1" {
				// 设置音频参数
				nodeInfoList = append(nodeInfoList, NodeInfo{
					NodeId:     param.NodeId,
					FieldName:  param.FieldName,
					FieldValue: audioResp.Data.FileName,
				})
				fmt.Printf("[音频] 设置到节点 %s: %s\n", param.NodeId, audioResp.Data.FileName)
			}
		} else {
			// 设置其他参数
			nodeInfoList = append(nodeInfoList, NodeInfo{
				NodeId:     param.NodeId,
				FieldName:  param.FieldName,
				FieldValue: param.FieldValue,
			})
		}
	}

	// 创建任务
	return CreateAdvancedTask(config.ID, nodeInfoList)
}

// MonitorTask 监控任务状态
func (we *WorkflowExecutor) MonitorTask(taskID string, onSuccess func(*TaskOutputResponse)) error {
	start := time.Now()
	for {
		statusResp, err := QueryTaskStatus(taskID)
		if err != nil {
			return fmt.Errorf("查询任务状态失败: %v", err)
		}

		elapsed := int(time.Since(start).Seconds())
		log.Printf("任务状态: %s (已等待 %d 秒)\n", statusResp.Data, elapsed)
		if statusResp.Data == "SUCCESS" || statusResp.Data == "FAILED" {
			totalElapsed := int(time.Since(start).Seconds())
			log.Printf("任务结束，最终状态: %s，总耗时: %d 秒\n", statusResp.Data, totalElapsed)
			if statusResp.Data == "SUCCESS" && onSuccess != nil {
				outputResp, err := QueryTaskOutputs(taskID)
				if err != nil {
					return fmt.Errorf("查询任务生成结果失败: %v", err)
				}
				onSuccess(outputResp)
			}
			break
		}

		time.Sleep(2 * time.Second)
	}
	return nil
}
