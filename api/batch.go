package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 创建结果保存目录
func createOutputDir() string {
	baseDir := "outputs"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}
	dateDir := filepath.Join(baseDir, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		log.Fatalf("创建日期目录失败: %v", err)
	}
	return dateDir
}

// 下载文件
func downloadFile(url, savePath string) error {
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}
	return nil
}

// 记录任务日志
func logTaskInfo(dir string, taskID string, outputs []TaskOutput) error {
	logFile := filepath.Join(dir, "task.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}
	defer file.Close()
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

// SaveTaskOutputs 保存任务输出结果到指定目录
func SaveTaskOutputs(outputDir, taskID string, outputs []TaskOutput, imageBaseName string) {
	for i, output := range outputs {
		fmt.Printf("[批量] - 文件URL: %s\n", output.FileUrl)
		fmt.Printf("[批量]   类型: %s\n", output.FileType)
		fmt.Printf("[批量]   节点ID: %s\n", output.NodeId)
		fmt.Printf("[批量]   任务耗时: %s\n", output.TaskCostTime)

		var fileName string
		if imageBaseName != "" {
			fileName = fmt.Sprintf("%s_%s_%d%s", imageBaseName, time.Now().Format("20060102_150405"), i, filepath.Ext(output.FileUrl))
		} else {
			fileName = fmt.Sprintf("%s_%s_%d%s", taskID, time.Now().Format("20060102_150405"), i, filepath.Ext(output.FileUrl))
		}
		savePath := filepath.Join(outputDir, fileName)
		if err := downloadFile(output.FileUrl, savePath); err != nil {
			fmt.Printf("[批量] 下载文件失败: %v\n", err)
			continue
		}
		fmt.Printf("[批量]   已保存到: %s\n", savePath)
	}
	// 记录任务日志
	if err := logTaskInfo(outputDir, taskID, outputs); err != nil {
		fmt.Printf("[批量] 记录任务日志失败: %v\n", err)
	}
}

// BatchProcessInputs 批量处理 inputs 目录下的图片文件
func BatchProcessInputs(workflowID string, concurrency int, executor *WorkflowExecutor) error {
	inputDir := "inputs"
	tmpDir := "tmp"
	outputDir := createOutputDir()

	// 创建 tmp 目录
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("创建 tmp 目录失败: %v", err)
	}

	// 获取所有文件
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("读取 inputs 目录失败: %v", err)
	}

	// 过滤文件
	var inputFiles []string
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".png" || ext == ".jpg" || ext == ".mp4" {
				inputFiles = append(inputFiles, filepath.Join(inputDir, file.Name()))
			}
		}
	}

	fmt.Printf("共获取到 %d 个文件：\n", len(inputFiles))
	for _, file := range inputFiles {
		fmt.Println("  -", file)
	}

	if len(inputFiles) == 0 {
		fmt.Println("inputs 目录下没有支持的文件（支持 .png、.jpg、.mp4）。")
		return nil
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, imgPath := range inputFiles {
		sem <- struct{}{}
		wg.Add(1)
		go func(img string) {
			defer wg.Done()
			fmt.Printf("[批量] 开始处理: %s\n", img)
			fmt.Printf("[批量] 上传图片: %s\n", img)
			imageBaseName := strings.TrimSuffix(filepath.Base(img), filepath.Ext(img))
			resp, err := executor.ExecuteWorkflowWithImage(workflowID, img)
			if err != nil {
				fmt.Printf("[批量] 处理失败: %s, 错误: %v\n", img, err)
				<-sem
				return
			} else if resp.Code != 0 || resp.Data.TaskId == "" {
				fmt.Printf("[批量] 任务创建失败: %s, code: %d, msg: %s\n", img, resp.Code, resp.Msg)
				<-sem
				return
			} else {
				fmt.Printf("[批量] 任务创建成功: %s, 任务ID: %s\n", img, resp.Data.TaskId)
				fmt.Printf("[批量] 等待任务完成: %s, 任务ID: %s\n", img, resp.Data.TaskId)
				err := executor.MonitorTask(resp.Data.TaskId, func(outputResp *TaskOutputResponse) {
					SaveTaskOutputs(outputDir, resp.Data.TaskId, outputResp.Data, imageBaseName)
				})
				if err != nil {
					fmt.Printf("[批量] 任务监控失败: %s, 错误: %v\n", img, err)
				}
				// 任务完成后立即移动文件
				dst := filepath.Join(tmpDir, filepath.Base(img))
				if err := os.Rename(img, dst); err != nil {
					fmt.Printf("[批量] 移动文件失败: %s -> %s, 错误: %v\n", img, dst, err)
				} else {
					fmt.Printf("[批量] 已移动到: %s\n", dst)
				}
			}
			<-sem
		}(imgPath)
	}

	wg.Wait()
	fmt.Println("批量处理完成。")
	return nil
}
