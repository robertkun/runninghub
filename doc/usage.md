# RunningHub 命令行工具使用指南

## 目录
- [基本命令](#基本命令)
- [工作流说明](#工作流说明)
- [常见问题](#常见问题)

## 基本命令

### 1. 列出所有可用工作流
```bash
go run main.go -list
```
此命令会显示所有已注册的工作流，包括：
- 工作流ID
- 工作流名称
- 工作流描述
- 节点配置信息

### 2. 单次处理
```bash
# 文本生成图片/视频
go run main.go -once -workflow <工作流ID>

# 图生图/图生视频
go run main.go -once -workflow <工作流ID> -image <图片路径>

# 视频+音频处理
go run main.go -once -workflow <工作流ID> -video <视频路径> -audio <音频路径>
```

### 3. 批量处理
```bash
# 批量处理图片
go run main.go -batchImg -workflow <工作流ID> [-concurrency N]

# 批量处理文本
go run main.go -batchText -workflow <工作流ID>
```

### 4. 任务管理
```bash
# 查询任务状态
go run main.go -task <任务ID>

# 取消任务
go run main.go -task <任务ID> -cancel
```

## 工作流说明

### 1. 图生视频工作流
- **WAN2.1 万相图生视频**
  - ID: `1930266544381792258`
  - 输入: 图片
  - 输出: 视频
  - 特点: 效果炸裂

- **泼水变装+换装**
  - ID: `1931186232649252865`
  - 输入: 图片
  - 输出: 视频
  - 特点: 通义万相，艾橘溪风格

### 2. 文生视频工作流
- **nunchaku-flux.1-dev+framePack**
  - ID: `1930520368543383553`
  - 输入: 文本提示词
  - 输出: 视频
  - 特点: 一键文生视频

### 3. 数字人工作流
- **数字人+口播**
  - ID: `1932095333768359938`
  - 输入: 视频 + 音频
  - 输出: 视频
  - 特点: 支持口播同步

### 4. VACE 14B 工作流
- **VACE 14B-图生视频**
  - ID: `1931521281978466306`
  - 输入: 图片
  - 输出: 视频
  - 特点: 支持自动提示词

## 常见问题

### 1. ApiKey 相关
- 确保在 `api/config.go` 中配置了正确的 ApiKey
- 不同 ApiKey 可能有不同的节点访问权限
- 如果遇到 `APIKEY_INVALID_NODE_INFO` 错误，请检查 ApiKey 权限

### 2. 文件格式支持
- 图片: `.png`, `.jpg`, `.jpeg`
- 视频: `.mp4`
- 音频: `.mp3`

### 3. 目录结构
```
.
├── inputs/           # 待处理文件目录
├── tmp/             # 已处理文件目录
├── outputs/         # 结果保存目录
│   └── YYYY-MM-DD/  # 按日期归档
└── doc/            # 文档目录
```

### 4. 错误处理
- 批量处理时，失败的文件会保留在 `inputs` 目录
- 成功的文件会被移动到 `tmp` 目录
- 所有任务日志保存在 `outputs/日期/task.log`

### 5. 性能优化
- 使用 `-concurrency` 参数控制并发数
- 默认并发数为 1
- 建议根据服务器性能调整并发数 