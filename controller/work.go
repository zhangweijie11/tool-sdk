package controller

import (
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gorm.io/gorm"
)

// TaskDeleteByTaskUUID 删除子任务和结果
func TaskDeleteByTaskUUID(taskUUID string) error {
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

// WorkDeleteByWorkUUID 删除总任务及其子任务和结果
func WorkDeleteByWorkUUID(workUUID string) error {
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

// WorkUpdateByWorkUUID 更新总任务状态
func WorkUpdateByWorkUUID(workUUID, column, newValue string) error {
	// 查询数据是否存在
	work, err := models.GetWorkByUUID(workUUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(schemas.RecordNotFoundErr)
	}
	// 更新总任务状态
	err = models.UpdateWork(work.UUID, column, newValue)
	if err != nil {
		return errors.New(schemas.RecordUpdateErr)
	}

	return nil
}

// WorkList 查看 Work 列表
func WorkList(schema *schemas.WorkListSchema) (*schemas.ListResponse, error) {
	var data = new(schemas.ListResponse)
	total, filter, query, err := models.WorkFilterQuery(&schema.Filter)
	if err != nil {
		return nil, err
	}
	data.TotalCount = total
	data.FilterCount = filter
	// 排序
	if len(schema.Order) == 0 {
		schema.Order = []string{"-update_time"}
	}
	for i := 0; i < len(schema.Order); i++ {
		switch schema.Order[i] {
		case "-update_time", "-source", "-id", "-create_time":
			query.Order(schema.Order[i][1:] + " desc")
		case "update_time", "source", "id", "create_time":
			query.Order(schema.Order[i])
		}
	}
	// 分页
	query = QueryPaging(query, schema.Page, schema.Size)
	var records []models.Work
	query.Find(&records)
	data.Records = records

	return data, nil
}

// WorkResultList 查看 Work 结果列表
func WorkResultList(schema *schemas.WorkResultListSchema) (*schemas.ListResponse, error) {
	var data = new(schemas.ListResponse)
	total, filter, query, err := models.WorkResultFilterQuery(&schema.Filter)
	if err != nil {
		return nil, err
	}
	data.TotalCount = total
	data.FilterCount = filter
	// 排序
	if len(schema.Order) == 0 {
		schema.Order = []string{"-update_time"}
	}
	for i := 0; i < len(schema.Order); i++ {
		switch schema.Order[i] {
		case "-update_time", "-source", "-id", "-create_time":
			query.Order(schema.Order[i][1:] + " desc")
		case "update_time", "source", "id", "create_time":
			query.Order(schema.Order[i])
		}
	}
	// 分页
	query = QueryPaging(query, schema.Page, schema.Size)
	var records []models.Result
	query.Find(&records)
	data.Records = records

	return data, nil
}
