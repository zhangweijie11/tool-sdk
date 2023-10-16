package models

import (
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"reflect"
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

// CreateResult 创建结果
func CreateResult(result *Result) error {
	if err := global.Db.Create(result).Error; err != nil {
		return err
	}

	return nil
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

func WorkResultFilterQuery(filter *schemas.WorkResultFilterSchema) (totalCount, filterCount int64, query *gorm.DB, err error) {
	query = global.Db.Model(&Result{})
	global.Db.Model(&Result{}).Count(&totalCount)
	// 判断搜索条件是否为空
	if reflect.DeepEqual(*filter, schemas.WorkFilterSchema{}) {
		filterCount = totalCount
	} else {
		if filter.WorkUUID != "" {
			query.Where("work_uuid = ?", filter.WorkUUID)
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
