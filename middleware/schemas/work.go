package schemas

// WorkCreateSchema 创建总任务参数
type WorkCreateSchema struct {
	WorkID   string                 `json:"task_id"`                     // 唯一标识
	Source   string                 `json:"source" binding:"required"`   // 任务来源
	Priority uint8                  `json:"priority" binding:"required"` // 任务优先级
	Params   map[string]interface{} `json:"params" binding:"required"`   // 任务参数
}

func (ws *WorkCreateSchema) ValidParams() error {
	return nil
}
