package domain

type PaginationResult[T any] struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"pageSize"`
	TotalPages   int   `json:"totalPages"`
	TotalRecords int64 `json:"totalRecords"`
	Data         []T   `json:"data"`
}
