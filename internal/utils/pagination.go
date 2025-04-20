package utils

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
	MaxPage         = 10000
)

// NormalizePageParams 统一分页参数校验与默认值
func NormalizePageParams(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if page > MaxPage {
		page = MaxPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}
