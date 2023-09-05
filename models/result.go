package models

import "gorm.io/datatypes"

type Result struct {
	ID       uint           `gorm:"column:id;primarykey" json:"id"`
	UUID     string         `gorm:"column:uuid;size:37;uniqueIndex;comment:结果uuid" json:"uuid"`
	WorkUUID string         `gorm:"column:work_uuid;size:37;index;comment:总任务uuid" json:"work_uuid"`
	TaskUUID string         `gorm:"column:work_uuid;size:37;index;comment:子任务uuid" json:"task_uuid"`
	Extra    datatypes.JSON `gorm:"column:extra;type:json,not null;default='{}';comment:全量数据" json:"extra"`
	ComplexBaseModel
}

func (Result) TableName() string {
	return "result"
}
