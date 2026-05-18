package pagination

type PageRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
