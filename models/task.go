package models

import "gorm.io/datatypes"

type Task struct {
	ID       uint           `gorm:"column:id;primarykey" json:"id"`
	UUID     string         `gorm:"column:uuid;size:37;uniqueIndex;comment:子任务uuid" json:"uuid"`
	WorkUUID string         `gorm:"column:work_uuid;size:37;index;not null;comment:总任务uuid" json:"work_uuid"`
	Params   datatypes.JSON `gorm:"column:params;type:json,not null;default='{}';comment:任务参数" json:"params"`
	//pending：待处理   doing：进行中   done：已完成   failed：失败   pause：暂停  cancelled：取消
	Status   string `gorm:"column:status;size:32;index;not null;default:'pending';comment:任务状态" json:"status"`
	Retry    uint8  `gorm:"column:retry;not null;default:0;comment:重试次数" json:"retry"`
	Priority uint8  `gorm:"column:priority;not null;default:0;comment:优先级" json:"priority"`
	Progress uint8  `gorm:"column:progress;not null;default:0;comment:进度" json:"progress"`
	ComplexBaseModel
}

func (Task) TableName() string {
	return "task"
}
