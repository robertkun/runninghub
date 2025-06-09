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

	// wan2.1_图生视频动作控制（跳舞工作流）
	GirlWorkflow = &WorkflowConfig{
		ID:          "1931185988360404994",
		Name:        "wan2.1_图生视频动作控制（跳舞工作流）",
		Description: "wan2.1_图生视频动作控制（跳舞工作流）",
		NodeConfigs: map[string]string{
			"18": "图片输入节点",
			"40": "文本提示词节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "34",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}

	// 泼水换装(图生视频)
	PoShuiWorkflow = &WorkflowConfig{
		ID:          "1931186232649252865",
		Name:        "泼水变装+换装=通义万相=艾橘溪",
		Description: "泼水变装+换装=通义万相=艾橘溪",
		NodeConfigs: map[string]string{
			"18": "图片输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "18",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}

	// VACE 14B-图生视频(可自动提示词)
	ZiZhuWorkflow = &WorkflowConfig{
		ID:          "1931350292812771329",
		Name:        "VACE 14B-图生视频(可自动提示词)",
		Description: "VACE 14B-图生视频(可自动提示词)",
		NodeConfigs: map[string]string{
			"210": "图片输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "210",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}
	
	// ATI字节最新轨迹驱动wan视频生成版
	ATiWorkflow = &WorkflowConfig{
		ID:          "1931384612306792449",
		Name:        "ATI字节最新轨迹驱动wan视频生成版)",
		Description: "ATI字节最新轨迹驱动wan视频生成版",
		NodeConfigs: map[string]string{
			"58": "图片输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "58",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}

	// FramePack-F1图生视频(帧率优化+时长优化+提示词优化)
	FramePackWorkflow = &WorkflowConfig{
		ID:          "1931386939079852033",
		Name:        "FramePack-F1图生视频(帧率优化+时长优化+提示词优化)",
		Description: "FramePack-F1图生视频(帧率优化+时长优化+提示词优化)",
		NodeConfigs: map[string]string{
			"2": "图片输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "2",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}

	// 运镜女孩(图生视频)
	OrbitWorkflow = &WorkflowConfig{
		ID:          "1931384038911574017",
		Name:        "FramePack-F1图生视频(帧率优化+时长优化+提示词优化)",
		Description: "FramePack-F1图生视频(帧率优化+时长优化+提示词优化)",
		NodeConfigs: map[string]string{
			"191": "图片输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "191",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}

	// FramePack-F1敏神最新算法图生视频
	FramePackF1Workflow = &WorkflowConfig{
		ID:          "1930526505527595010",
		Name:        "FramePack-F1敏神最新算法图生视频",
		Description: "FramePack-F1敏神最新算法图生视频",
		NodeConfigs: map[string]string{
			"2": "图片输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "2",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}

	// VACE 14B-图生视频(可自动提示词)
	VACE14BWorkflow = &WorkflowConfig{
		ID:          "1931521281978466306",
		Name:        "VACE 14B-图生视频(可自动提示词)",
		Description: "VACE 14B-图生视频(可自动提示词)",
		NodeConfigs: map[string]string{
			"210": "图片输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "210",
				FieldName:  "image",
				FieldValue: "", // 图片路径会在执行时设置
				IsImage:    true,
			},
		},
	}

	// 数字人+口播
	ShuZiRenWorkflow = &WorkflowConfig{
		ID:          "1927694627866820610",
		Name:        "数字人+口播",
		Description: "数字人+口播",
		NodeConfigs: map[string]string{
			"1": "视频输入节点",
		},
		Params: []NodeParam{
			{
				NodeId:     "1",
				FieldName:  "video",
				FieldValue: "", // 视频路径会在执行时设置
				IsImage:    true,  // 改回 true，因为我们需要使用图片上传接口
			},
		},
	}
)
