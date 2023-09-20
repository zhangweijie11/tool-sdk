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
