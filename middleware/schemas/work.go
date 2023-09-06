package schemas

// WorkCreateSchema 创建总任务参数
type WorkCreateSchema struct {
	WorkUUID string                 `json:"work_uuid"`                                // 唯一标识
	Source   string                 `json:"source" binding:"required"`                // 任务来源
	Priority uint8                  `json:"priority" binding:"required,max=9,min=1" ` // 任务优先级
	Params   map[string]interface{} `json:"params" binding:"required"`                // 任务参数
}

// WorkDeleteSchema 删除总任务参数
type WorkDeleteSchema struct {
	WorkUUID string `json:"work_uuid" binding:"required"` // 唯一标识
}

// WorkGetInfoSchema 获取总任务参数
type WorkGetInfoSchema struct {
	WorkUUID string `json:"work_uuid" binding:"required"` // 唯一标识
}
