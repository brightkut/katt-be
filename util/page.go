package util

import "gorm.io/gorm"

type PaginationResult[T any] struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"pageSize"`
	TotalPages   int   `json:"totalPages"`
	TotalRecords int64 `json:"totalRecords"`
	Data         []T   `json:"data"`
}

func Paginate[T any](db *gorm.DB, model T, page, pageSize int) (PaginationResult[T], error) {
	var result PaginationResult[T]
	var data []T
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	db.Model(&model).Count(&total)

	offset := (page - 1) * pageSize
	err := db.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&data).Error
	if err != nil {
		return result, err
	}

	result = PaginationResult[T]{
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   int((total + int64(pageSize) - 1) / int64(pageSize)),
		TotalRecords: total,
		Data:         data,
	}

	return result, nil
}
