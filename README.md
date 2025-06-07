# RunningHub 批量与单次任务处理工具

## 项目简介

本项目是基于 Go 语言开发的 RunningHub 工作流自动化处理工具，支持图片/视频等任务的单次处理与批量处理，自动保存结果，支持任务状态查询、并发控制等功能，适用于 AI 生成类场景。

---

## 主要功能

- **单次处理**：支持通过命令行一次性处理一张图片，自动等待任务完成并保存结果。
- **批量处理**：支持批量处理 `inputs` 目录下的所有图片，支持自定义并发数，任务全流程（上传、创建、执行完成）后再处理下一张。
- **任务状态查询**：可通过任务ID实时查询任务执行状态和结果。
- **并发控制**：批量处理时可通过 `-concurrency` 参数设置并发数，默认1。
- **结果保存**：所有生成结果自动保存到 `outputs/日期/` 目录下，文件名包含任务ID、时间戳、图片名等信息。
- **失败重试**：批量处理时，任务失败的图片不会被移动，便于后续排查和重试。
- **日志记录**：每次任务执行会在输出目录下生成 `task.log`，记录详细任务信息。

---

## 依赖与环境

- Go 1.18 及以上
- 需联网访问 RunningHub API
- 需在 `api/config.go` 中配置有效的 `ApiKey`

---

## 目录结构

```
.
├── main.go           # 主程序入口
├── api/
│   ├── workflow.go   # 工作流配置与管理
│   ├── executor.go   # 工作流执行与任务监控
│   ├── task.go       # 任务API调用
│   ├── upload.go     # 图片上传API
│   ├── batch.go      # 批量处理逻辑
│   └── config.go     # API Key 配置
├── inputs/           # 批量处理时待处理图片目录
├── tmp/              # 批量处理后已处理图片目录
├── outputs/          # 结果保存目录，按日期归档
│   └── YYYY-MM-DD/   # 每天的结果子目录
│       ├── *.png     # 生成图片/视频等
│       └── task.log  # 任务日志
└── README.md         # 使用说明
```

---

## 使用方法

### 1. 列出所有可用工作流
```bash
go run main.go -list
```

### 2. 单次处理（需加 -once）
```bash
go run main.go -once -workflow <工作流ID> [-image <图片路径>]
```
- 例：
  - 文本生成图片：`go run main.go -once -workflow 1930266544381792258`
  - 图生图/图生视频：`go run main.go -once -workflow 1930266544381792258 -image test.png`

### 3. 批量处理
```bash
go run main.go -batch -workflow <工作流ID> [-concurrency N]
```
- 默认并发为1，支持自定义并发数
- 例：
  - 串行：`go run main.go -batch -workflow 1930266544381792258`
  - 并发3：`go run main.go -batch -workflow 1930266544381792258 -concurrency 3`
- 处理完成后，成功的图片会被移动到 `tmp/`，失败的图片保留在 `inputs/`

### 4. 查询任务状态
```bash
go run main.go -task <任务ID>
```

### 5. 取消任务
```bash
go run main.go -task <任务ID> -cancel
```

---

## 注意事项

- 请确保 `inputs/` 目录下有待处理图片，支持 `.png`、`.jpg`、`.jpeg` 格式
- 需在 `api/config.go` 配置有效的 `ApiKey`
- 工作流配置需在 `api/workflow.go` 注册
- 结果文件和日志自动保存到 `outputs/日期/` 目录
- 批量处理时，只有任务创建并执行完成的图片才会被移动到 `tmp/`
- 失败的图片不会被移动，便于后续重试

---

## 联系方式
如有问题或建议，请联系开发者。 