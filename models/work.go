package models

import (
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"reflect"
)

type Work struct {
	ID             uint           `gorm:"column:id;primarykey" json:"id"`
	UUID           string         `gorm:"column:uuid;size:128;uniqueIndex;comment:总任务uuid" json:"uuid"`
	Params         datatypes.JSON `gorm:"column:params;type:json,not null;default='{}';not null;comment:总任务参数" json:"params"` //pending：待处理   doing：进行中   done：已完成   failed：失败   pause：暂停  cancelled：取消
	Status         string         `gorm:"column:status;size:32;index;not null;default:'pending';comment:总任务状态" json:"status"`
	Source         string         `gorm:"column:source;size:64;index;not null;default:'';comment:总任务来源" json:"source"`
	CallbackUrl    string         `gorm:"column:callback_url;size:128;not null;default:'';comment:回调地址" json:"callback_url"`
	CallbackType   string         `gorm:"column:callback_type;size:64;not null;default:'';comment:回调方式" json:"callback_type"`
	CallbackStatus string         `gorm:"column:callback_status;size:64;not null;default:'';comment:回调状态" json:"callback_status"`
	WorkType       string         `gorm:"column:work_type;size:64;not null;default:'';comment:任务类型" json:"work_type"`
	Retry          uint8          `gorm:"column:retry;not null;default:0;comment:重试次数" json:"retry"`
	Priority       uint8          `gorm:"column:priority;not null;default:0;comment:优先级" json:"priority"`
	Progress       uint8          `gorm:"column:progress;not null;default:0;comment:进度" json:"progress"`
	ProgressType   string         `gorm:"column:progress_type;not null;default:'';comment:进度推送方式" json:"progress_type"`
	ProgressUrl    string         `gorm:"column:progress_url;not null;default:'';comment:进度推送地址" json:"progress_url"`
	ComplexBaseModel
}

func (Work) TableName() string {
	return "work"
}

// CreateWok 创建 work 总任务
func CreateWok(work *Work) error {
	if err := global.Db.Create(work).Error; err != nil {
		return err
	}

	return nil
}

// DeleteWorkByWorkUUID 根据总任务唯一标识删除数据
func DeleteWorkByWorkUUID(workUUID string) error {
	if err := global.Db.Where("uuid = ?", workUUID).Delete(&Work{}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteWorkByWorkUUIDs 根据总任务唯一标识批量删除数据
func DeleteWorkByWorkUUIDs(workUUIDs []string) error {
	if err := global.Db.Model(&Work{}).Where("uuid in ?", workUUIDs).Delete(&Work{}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateWork 创建 work 总任务
func UpdateWork(workUUID, column, newValue string) error {
	if err := global.Db.Model(&Work{}).Where("uuid = ?", workUUID).Update(column, newValue).Error; err != nil {
		return err
	}

	return nil
}

// GetWorkByUUID 根据总任务唯一标识查询数据
func GetWorkByUUID(workUUID string) (Work, error) {
	var work Work
	if err := global.Db.Where("uuid = ?", workUUID).First(&work).Error; err != nil {
		return work, err
	}

	return work, nil
}

// GetWorkOrderCreateTime 根据创建时间排序获取待执行的任务
func GetWorkOrderCreateTime() (Work, error) {
	var work Work
	if err := global.Db.Where("status in (?)", []string{global.WorkStatusPending, global.WorkStatusRestart}).Order("create_time asc").Find(&work).Error; err != nil {
		return work, err
	}

	return work, nil
}

// UpdateWorkDoingToPending 将状态为 doing 的任务变更为 pending
func UpdateWorkDoingToPending() error {
	if global.Config.Server.ServerName != "" {
		if err := global.Db.Model(&Work{}).Where("status = ? AND work_type = ?", global.WorkStatusDoing, global.Config.Server.ServerName).Update("status", global.WorkStatusPending).Error; err != nil {
			return err
		}
	} else {
		if err := global.Db.Model(&Work{}).Where("status = ?", global.WorkStatusDoing).Update("status", global.WorkStatusPending).Error; err != nil {
			return err
		}
	}

	return nil
}

func WorkFilterQuery(filter *schemas.WorkFilterSchema) (totalCount, filterCount int64, query *gorm.DB, err error) {
	query = global.Db.Model(&Work{})
	global.Db.Model(&Work{}).Count(&totalCount)
	// 判断搜索条件是否为空
	if reflect.DeepEqual(*filter, schemas.WorkFilterSchema{}) {
		filterCount = totalCount
	} else {
		if filter.WorkUUID != "" {
			query.Where("uuid = ?", filter.WorkUUID)
		}
		if filter.WorkType != "" {
			query.Where("work_type = ?", filter.WorkType)
		}
		if filter.WorkStatus != "" {
			query.Where("status = ?", filter.WorkStatus)
		}
		if filter.WorkSource != "" {
			query.Where("source = ?", filter.WorkSource)
		}
		if filter.WorkPriority != 0 {
			query.Where("priority = ?", filter.WorkPriority)
		}
		// 创建时间
		if filter.CreateTime != nil {
			times, err := schemas.TimeRangeValidator(filter.CreateTime)
			if err != nil {
				return totalCount, filterCount, query, errors.New(schemas.TimeCreateErr)
			}
			query.Where("create_time BETWEEN ? AND ?", times[0], times[1])
		}

		// 更新时间
		if filter.UpdateTime != nil {
			times, err := schemas.TimeRangeValidator(filter.UpdateTime)
			if err != nil {
				return totalCount, filterCount, query, errors.New(schemas.TimeUpdateErr)
			}
			query.Where("update_time BETWEEN ? AND ?", times[0], times[1])
		}
	}
	query.Count(&filterCount)

	return totalCount, filterCount, query, err
}
