package api

// NodeParam 节点参数配置
type NodeParam struct {
	NodeId     string      `json:"nodeId"`     // 节点ID
	FieldName  string      `json:"fieldName"`  // 字段名
	FieldValue interface{} `json:"fieldValue"` // 字段值
	IsImage    bool        `json:"isImage"`    // 是否为图片输入节点
}

// WorkflowConfig 工作流配置
type WorkflowConfig struct {
	ID          string            `json:"id"`          // 工作流ID
	Name        string            `json:"name"`        // 工作流名称
	Description string            `json:"description"` // 工作流描述
	NodeConfigs map[string]string `json:"nodeConfigs"` // 节点配置，key为节点ID，value为节点描述
	Params      []NodeParam       `json:"params"`      // 固定参数配置
}

// WorkflowManager 工作流管理器
type WorkflowManager struct {
	workflows map[string]*WorkflowConfig
}

// NewWorkflowManager 创建工作流管理器
func NewWorkflowManager() *WorkflowManager {
	return &WorkflowManager{
		workflows: make(map[string]*WorkflowConfig),
	}
}

// RegisterWorkflow 注册工作流
func (wm *WorkflowManager) RegisterWorkflow(config *WorkflowConfig) {
	wm.workflows[config.ID] = config
}

// GetWorkflow 获取工作流配置
func (wm *WorkflowManager) GetWorkflow(id string) (*WorkflowConfig, bool) {
	config, exists := wm.workflows[id]
	return config, exists
}

// ListWorkflows 列出所有工作流
func (wm *WorkflowManager) ListWorkflows() []*WorkflowConfig {
	configs := make([]*WorkflowConfig, 0, len(wm.workflows))
	for _, config := range wm.workflows {
		configs = append(configs, config)
	}
	return configs
}

// 预定义的工作流配置
var (
	// 图生视频(WAN2.1 相)
	CatWorkflow = &WorkflowConfig{
		ID:          "1930266544381792258",
		Name:        "WAN2.1 万相图生视频，效果炸裂",
		Description: "WAN2.1 万相图生视频，效果炸裂",
		NodeConfigs: map[string]string{
			"18": "图片输入节点",
			"40": "文本提示词节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "18",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
			{
				NodeId:     "40",
				FieldName:  "text",
				FieldValue: "",
				IsImage:    false,
			},
		},
	}

	// 文生视频(nunchaku-flux.1-dev+framePack 一键文生视频)
	DogWorkflow = &WorkflowConfig{
		ID:          "1930520368543383553",
		Name:        "nunchaku-flux.1-dev+framePack 一键文生视频",
		Description: "nunchaku-flux.1-dev+framePack 一键文生视频",
		NodeConfigs: map[string]string{
			"6": "文本提示词节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "6",
				FieldName:  "text",
				FieldValue: "Realistic style, The Little Girl with Matchsticks, ",
				IsImage:    false,
			},
		},
	}

	// 可以继续添加更多工作流配置...
)
