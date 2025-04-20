package schema

// PageRequest 通用分页请求结构体（可嵌入到具体 ListXXXRequest）
type PageRequest struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	OrderBy  string `json:"order_by" query:"order_by"` // 排序字段
	Desc     bool   `json:"desc" query:"desc"`         // 是否降序
}

// PageResult 通用分页响应结构体（可嵌入到具体 ListXXXResponse）
type PageResult[T any] struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	Items      []T   `json:"items"`
}
