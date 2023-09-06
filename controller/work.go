package controller

import (
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gorm.io/gorm"
)

// DeleteTaskByTaskUUID 删除子任务和结果
func DeleteTaskByTaskUUID(taskUUID string) error {
	// 查询数据是否存在
	task, err := models.GetTaskByTaskUUID(taskUUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(schemas.RecordNotFoundErr)
	}

	// 删除子任务
	err = models.DeleteTaskByTaskUUID(task.UUID)
	if err != nil {
		return errors.New(schemas.RecordDeleteErr)
	}

	// 删除结果
	err = models.DeleteResultByTaskUUID(task.UUID)
	if err != nil {
		return errors.New(schemas.RecordDeleteErr)
	}

	return nil
}

// DeleteWorkByWorkUUID 删除总任务及其子任务和结果
func DeleteWorkByWorkUUID(workUUID string) error {
	// 查询数据是否存在
	work, err := models.GetWorkByUUID(workUUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(schemas.RecordNotFoundErr)
	}
	// 删除总任务
	err = models.DeleteWorkByWorkUUID(work.UUID)
	if err != nil {
		return errors.New(schemas.RecordDeleteErr)
	}

	// 删除子任务
	err = models.DeleteTaskByWorkUUID(work.UUID)
	if err != nil {
		return errors.New(schemas.RecordDeleteErr)
	}

	// 删除结果
	err = models.DeleteResultByWorkUUID(work.UUID)
	if err != nil {
		return errors.New(schemas.RecordDeleteErr)
	}

	return nil
}

// UpdateWorkByWorkUUID 更新总任务状态
func UpdateWorkByWorkUUID(workUUID, column, newValue string) error {
	// 查询数据是否存在
	work, err := models.GetWorkByUUID(workUUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(schemas.RecordNotFoundErr)
	}
	// 更新总任务状态
	err = models.UpdateWork(work.UUID, column, newValue)
	if err != nil {
		return errors.New(schemas.RecordDeleteErr)
	}

	return nil
}
