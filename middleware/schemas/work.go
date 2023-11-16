package schemas

// WorkCreateSchema 创建总任务参数
type WorkCreateSchema struct {
	WorkUUID     string                 `json:"work_uuid"`                                // 唯一标识
	Source       string                 `json:"source" binding:"required"`                // 任务来源
	Priority     uint8                  `json:"priority" binding:"required,max=9,min=1" ` // 任务优先级
	Params       map[string]interface{} `json:"params" binding:"required"`                // 任务参数
	CallbackUrl  string                 `json:"callback_url"`                             // 回调地址
	CallbackType string                 `json:"callback_type"`                            // 回调方式
	ProgressType string                 `json:"progress_type"`                            // 进度推送方式
	ProgressUrl  string                 `json:"progress_url"`                             // 进度推送地址
	WorkType     string                 `json:"work_type"`                                // 任务类型
}

// WorkDeleteSchema 删除总任务参数
type WorkDeleteSchema struct {
	WorkUUID string `json:"work_uuid" binding:"required"` // 唯一标识
}

// WorkGetInfoSchema 获取总任务参数
type WorkGetInfoSchema struct {
	WorkUUID string `json:"work_uuid" binding:"required"` // 唯一标识
}

// WorkPauseSchema 暂停总任务参数
type WorkPauseSchema struct {
	WorkUUID string `json:"work_uuid" binding:"required"` // 唯一标识
}

// WorkStopSchema 停止总任务参数
type WorkStopSchema struct {
	WorkUUID string `json:"work_uuid" binding:"required"` // 唯一标识
}

// WorkRestartSchema 重启总任务参数
type WorkRestartSchema struct {
	WorkUUID string `json:"work_uuid" binding:"required"` // 唯一标识
}

type WorkFilterSchema struct {
	WorkUUID     string `json:"work_uuid"`
	WorkType     string `json:"work_type"`
	WorkStatus   string `json:"work_status"`
	WorkSource   string `json:"work_source"`
	WorkPriority uint8  `json:"work_priority" binding:"min=0,max=9"`
	// 范围匹配
	CreateTime []string `json:"create_time" binding:"omitempty,len=2"` // 创建时间
	UpdateTime []string `json:"update_time" binding:"omitempty,len=2"` // 更新时间
}

// WorkListSchema 任务列表
type WorkListSchema struct {
	Page   int              `json:"page" binding:"min=0"`              // 请求页码 default(1) min(1)
	Size   int              `json:"size" binding:"min=0,max=100"`      // 请求数量 default(10)，range(1-100)
	Order  []string         `json:"order"  binding:"omitempty,unique"` // 排序
	Filter WorkFilterSchema `json:"filter"`                            // 筛选条件
}

type WorkResultFilterSchema struct {
	WorkUUID string `json:"work_uuid"`
	// 范围匹配
	CreateTime []string `json:"create_time" binding:"omitempty,len=2"` // 创建时间
	UpdateTime []string `json:"update_time" binding:"omitempty,len=2"` // 更新时间
}

// WorkResultListSchema 任务结果列表
type WorkResultListSchema struct {
	Page   int                    `json:"page" binding:"min=0"`              // 请求页码 default(1) min(1)
	Size   int                    `json:"size" binding:"min=0,max=100"`      // 请求数量 default(10)，range(1-100)
	Order  []string               `json:"order"  binding:"omitempty,unique"` // 排序
	Filter WorkResultFilterSchema `json:"filter"`                            // 筛选条件
}
