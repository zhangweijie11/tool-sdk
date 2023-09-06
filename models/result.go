package models

import (
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gorm.io/datatypes"
)

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

// DeleteResultByWorkUUID 根据总任务唯一标识删除数据
func DeleteResultByWorkUUID(workUUID string) error {
	if err := global.Db.Where("work_uuid = ?", workUUID).Delete(&Result{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteResultByTaskUUID 根据子任务唯一标识删除数据
func DeleteResultByTaskUUID(taskUUID string) error {
	if err := global.Db.Where("uuid = ?", taskUUID).Delete(&Result{}).Error; err != nil {
		return err
	}
	return nil
}
