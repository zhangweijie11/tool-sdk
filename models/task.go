package models

import (
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gorm.io/datatypes"
)

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

// CreateTask 创建 task 子任务
func CreateTask(task *Task) error {
	if err := global.Db.Create(task).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTaskByWorkUUID 根据总任务唯一标识删除数据
func DeleteTaskByWorkUUID(workUUID string) error {
	if err := global.Db.Where("work_uuid = ?", workUUID).Delete(&Task{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTaskByTaskUUID 根据子任务唯一标识删除数据
func DeleteTaskByTaskUUID(taskUUID string) error {
	if err := global.Db.Where("uuid = ?", taskUUID).Delete(&Task{}).Error; err != nil {
		return err
	}
	return nil
}

// GetTaskByTaskUUID 根据子任务唯一标识查询数据
func GetTaskByTaskUUID(taskUUID string) (Task, error) {
	var task Task
	if err := global.Db.Where("uuid = ?", taskUUID).First(&task).Error; err != nil {
		return task, err
	}

	return task, nil
}
