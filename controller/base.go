package controller

import "gorm.io/gorm"

func QueryPaging(query *gorm.DB, page, size int) *gorm.DB {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return query.Offset((page - 1) * size).Limit(size)
}
