package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"strconv"

	"runninghub/api"
)

// 创建结果保存目录
func createOutputDir() string {
	// 创建基础目录
	baseDir := "outputs"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 创建日期子目录
	dateDir := filepath.Join(baseDir, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		log.Fatalf("创建日期目录失败: %v", err)
	}

	return dateDir
}

// 下载文件
func downloadFile(url, savePath string) error {
	// 创建文件
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()

	// 保存文件
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// 记录任务日志
func logTaskInfo(dir string, taskID string, outputs []api.TaskOutput) error {
	logFile := filepath.Join(dir, "task.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}
	defer file.Close()

	// 写入任务信息
	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("任务ID: %s", taskID)
	logger.Printf("执行时间: %s", time.Now().Format("2006-01-02 15:04:05"))
	logger.Println("生成结果:")
	for _, output := range outputs {
		logger.Printf("- 文件URL: %s", output.FileUrl)
		logger.Printf("  类型: %s", output.FileType)
		logger.Printf("  节点ID: %s", output.NodeId)
		logger.Printf("  任务耗时: %s", output.TaskCostTime)
	}
	logger.Println("----------------------------------------")

	return nil
}

func BatchProcessText(workflowID string, executor *api.WorkflowExecutor) error {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 构建文件的绝对路径
	filePath := filepath.Join(wd, "doc", "book.txt")
	fmt.Printf("尝试读取文件: %s\n", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}
	fmt.Printf("文件大小: %d 字节\n", fileInfo.Size())

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	if len(content) == 0 {
		return fmt.Errorf("文件内容为空")
	}

	fmt.Printf("读取到的内容长度: %d 字节\n", len(content))

	// 统一换行符为 \n，然后按行分割
	contentStr := strings.ReplaceAll(string(content), "\r\n", "\n")
	contentStr = strings.ReplaceAll(contentStr, "\r", "\n")
	paragraphs := strings.Split(contentStr, "\n")

	outputDir := createOutputDir()
	fmt.Println("outputDir: ", outputDir)
	fmt.Println("paragraphs: ", len(paragraphs))

	// 打印每个段落的内容（用于调试）
	for i, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		fmt.Printf("段落 %d 长度: %d\n", i+1, len(p))
	}

	for idx, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		fmt.Printf("[批量文本] 开始处理第 %d 段: %s\n", idx+1, para)

		// 执行工作流，使用当前段落作为文本参数
		resp, err := executor.ExecuteWorkflowWithText(workflowID, para)
		if err != nil {
			fmt.Printf("[批量文本] 处理失败: %v\n", err)
			continue
		}
		if resp.Code != 0 || resp.Data.TaskId == "" {
			fmt.Printf("[批量文本] 任务创建失败: code=%d, msg=%s\n", resp.Code, resp.Msg)
			continue
		}
		fmt.Printf("[批量文本] 任务创建成功: 任务ID: %s\n", resp.Data.TaskId)
		fmt.Printf("[批量文本] 等待任务完成: 任务ID: %s\n", resp.Data.TaskId)

		err = executor.MonitorTask(resp.Data.TaskId, func(outputResp *api.TaskOutputResponse) {
			fmt.Printf("[批量文本] 任务完成: 任务ID: %s\n", resp.Data.TaskId)
			// 保存输出
			imageBaseName := fmt.Sprintf("text_%d", idx+1)
			api.SaveTaskOutputs(outputDir, resp.Data.TaskId, outputResp.Data, imageBaseName)
		})
		if err != nil {
			fmt.Printf("[批量文本] 任务监控失败: %v\n", err)
		}
		// 顺序执行，等待当前任务完成后再处理下一个
	}
	fmt.Println("批量文本处理完成。")
	return nil
}

func main() {
	// 启动时获取并打印账户信息
	apiKey := api.GetApiKey()
	if apiKey == "" {
		fmt.Println("[警告] 未设置 API Key，无法获取账户信息。")
	} else {
		status, err := api.GetAccountStatus(apiKey)
		if err != nil {
			fmt.Printf("[账户信息] 获取失败: %v\n", err)
		} else if status.Code != 0 {
			fmt.Printf("[账户信息] 获取失败: %s\n", status.Msg)
		} else {
			remainCoins := status.Data.RemainCoins
			currentTaskCounts := status.Data.CurrentTaskCounts
			// 尝试转换为 int
			if coins, err := strconv.Atoi(remainCoins); err == nil {
				remainCoins = fmt.Sprintf("%d", coins)
			}
			if tasks, err := strconv.Atoi(currentTaskCounts); err == nil {
				currentTaskCounts = fmt.Sprintf("%d", tasks)
			}
			fmt.Printf("[账户信息] 剩余金币: %s，当前任务数: %s\n", remainCoins, currentTaskCounts)
		}
	}
	time.Sleep(1 * time.Second)

	// 定义命令行参数
	taskID := flag.String("task", "", "要查询的任务ID")
	cancel := flag.Bool("cancel", false, "是否取消任务")
	workflowID := flag.String("workflow", "", "要执行的工作流ID")
	imagePath := flag.String("image", "", "要上传的图片路径")
	videoPath := flag.String("video", "", "要上传的视频路径")
	audioPath := flag.String("audio", "", "要上传的音频路径")
	list := flag.Bool("list", false, "列出所有可用的工作流")
	batchImg := flag.Bool("batchImg", false, "批量处理 inputs 目录下的图片")
	batchText := flag.Bool("batchText", false, "批量处理 inputs 目录下的图片")
	once := flag.Bool("once", false, "批量处理 inputs 目录下的图片")
	concurrency := flag.Int("concurrency", 1, "并发数量")
	flag.Parse()

	// 创建工作流管理器
	manager := api.NewWorkflowManager()

	// 注册工作流
	manager.RegisterWorkflow(api.CatWorkflow)
	manager.RegisterWorkflow(api.DogWorkflow)
	manager.RegisterWorkflow(api.GirlWorkflow)
	manager.RegisterWorkflow(api.PoShuiWorkflow)
	manager.RegisterWorkflow(api.ZiZhuWorkflow)
	manager.RegisterWorkflow(api.ATiWorkflow)
	manager.RegisterWorkflow(api.FramePackWorkflow)
	manager.RegisterWorkflow(api.FramePackF1Workflow)
	manager.RegisterWorkflow(api.OrbitWorkflow)
	manager.RegisterWorkflow(api.VACE14BWorkflow)
	manager.RegisterWorkflow(api.ShuZiRenWorkflow)
	manager.RegisterWorkflow(api.VACE14BWorkflow2)

	// 在这里注册更多工作流...

	// 创建工作流执行器
	executor := api.NewWorkflowExecutor(manager)

	switch {
	case *batchImg: 
		if *workflowID == "" {
			log.Fatalf("批量处理时必须指定 -workflow <工作流ID>")
		}
		err := api.BatchProcessInputs(*workflowID, *concurrency, executor)
		if err != nil {
			log.Fatalf("批量处理失败: %v", err)
		}
		return
	case *batchText:
		if *workflowID == "" {
			log.Fatalf("批量处理时必须指定 -workflow <工作流ID>")
		}
		err := BatchProcessText(*workflowID, executor)
		if err != nil {
			log.Fatalf("批量文本处理失败: %v", err)
		}
		return

	case *once:
		if *workflowID == "" {
			log.Fatalf("单次处理时必须指定 -workflow <工作流ID>")
		}
		var resp *api.TaskCreateResponse
		var err error

		var imageBaseName string
		if *videoPath != "" && *audioPath != "" {
			// 执行带视频和音频的工作流
			resp, err = executor.ExecuteWorkflowWithVideoAndAudio(*workflowID, *videoPath, *audioPath)
			// 使用视频文件名作为基础名
			base := filepath.Base(*videoPath)
			imageBaseName = strings.TrimSuffix(base, filepath.Ext(base))
		} else if *imagePath != "" {
			// 获取图片基础名（不含扩展名）
			base := filepath.Base(*imagePath)
			imageBaseName = strings.TrimSuffix(base, filepath.Ext(base))
			// 执行带图片的工作流
			resp, err = executor.ExecuteWorkflowWithImage(*workflowID, *imagePath)
		} else {
			// 执行普通工作流
			resp, err = executor.ExecuteWorkflow(*workflowID)
		}

		if err != nil {
			log.Fatalf("执行工作流失败: %v", err)
		}

		// 检查任务创建是否成功
		if resp.Code != 0 || resp.Data.TaskId == "" {
			log.Fatalf("任务创建失败！code: %d, msg: %s", resp.Code, resp.Msg)
		}

		fmt.Printf("任务创建成功！任务ID: %s\n", resp.Data.TaskId)
		fmt.Println("正在等待任务完成...")

		// 创建输出目录
		outputDir := createOutputDir()

		// 自动监控任务状态并显示结果
		err = executor.MonitorTask(resp.Data.TaskId, func(outputResp *api.TaskOutputResponse) {
			fmt.Println("\n任务执行成功！")
			fmt.Println("生成结果:")
			timestamp := time.Now().Format("20060102_150405")
			for i, output := range outputResp.Data {
				fmt.Printf("- 文件URL: %s\n", output.FileUrl)
				fmt.Printf("  类型: %s\n", output.FileType)
				fmt.Printf("  节点ID: %s\n", output.NodeId)
				fmt.Printf("  任务耗时: %s\n", output.TaskCostTime)

				// 判断是否为视频类型，若是则用图片名命名
				var fileName string
				if imageBaseName != "" {
					fileName = fmt.Sprintf("%s_%s_%d%s", imageBaseName, timestamp, i, filepath.Ext(output.FileUrl))
				} else {
					fileName = fmt.Sprintf("%s_%s_%d%s", resp.Data.TaskId, timestamp, i, filepath.Ext(output.FileUrl))
				}
				savePath := filepath.Join(outputDir, fileName)
				if err := downloadFile(output.FileUrl, savePath); err != nil {
					log.Printf("下载文件失败: %v", err)
					continue
				}
				fmt.Printf("  已保存到: %s\n", savePath)
			}

			// 记录任务日志
			if err := logTaskInfo(outputDir, resp.Data.TaskId, outputResp.Data); err != nil {
				log.Printf("记录任务日志失败: %v", err)
			}
		})
		if err != nil {
			log.Fatalf("监控任务失败: %v", err)
		}
		return

	case *taskID != "":
		if *cancel {
			// 取消任务
			resp, err := api.CancelTask(*taskID)
			if err != nil {
				log.Fatalf("取消任务失败: %v", err)
			}
			fmt.Printf("取消任务响应: %+v\n", resp)
			return
		}

		// 监控任务状态
		err := executor.MonitorTask(*taskID, func(outputResp *api.TaskOutputResponse) {
			fmt.Println("\n任务执行成功！")
			fmt.Println("生成结果:")
			for _, output := range outputResp.Data {
				fmt.Printf("- 文件URL: %s\n", output.FileUrl)
				fmt.Printf("  类型: %s\n", output.FileType)
				fmt.Printf("  节点ID: %s\n", output.NodeId)
				fmt.Printf("  任务耗时: %s\n", output.TaskCostTime)
			}
		})
		if err != nil {
			log.Fatalf("监控任务失败: %v", err)
		}
		return

	case *list:
		workflows := manager.ListWorkflows()
		fmt.Println("可用的工作流:")
		for _, wf := range workflows {
			fmt.Printf("\n工作流ID: %s\n", wf.ID)
			fmt.Printf("名称: %s\n", wf.Name)
			fmt.Printf("描述: %s\n", wf.Description)
			fmt.Println("节点配置:")
			for _, param := range wf.Params {
				if param.IsImage {
					fmt.Printf("  - 图片输入节点: %s\n", param.NodeId)
				} else {
					fmt.Printf("  - 节点: %s, 字段: %s, 值: %s\n", param.NodeId, param.FieldName, param.FieldValue)
				}
			}
		}
		return

	default:
		fmt.Println("使用方法:")
		fmt.Println("1. 列出所有工作流:")
		fmt.Println("   go run main.go -list")
		fmt.Println("\n2. 执行工作流:")
		fmt.Println("   go run main.go -workflow <工作流ID>")
		fmt.Println("\n3. 执行带图片的工作流:")
		fmt.Println("   go run main.go -workflow <工作流ID> -image <图片路径>")
		fmt.Println("\n4. 执行对口型工作流:")
		fmt.Println("   go run main.go -workflow <工作流ID> -video <视频路径> -audio <音频路径>")
		fmt.Println("\n5. 查询任务状态:")
		fmt.Println("   go run main.go -task <任务ID>")
		fmt.Println("\n6. 取消任务:")
		fmt.Println("   go run main.go -task <任务ID> -cancel")
		fmt.Println("\n7. 批量处理:")
		fmt.Println("   go run main.go -batch -workflow <工作流ID> [-concurrency N]")
	}
}
