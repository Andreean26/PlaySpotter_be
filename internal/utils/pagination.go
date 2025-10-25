package utils

type PaginationParams struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type PaginationMeta struct {
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PageCount int   `json:"page_count"`
	Limit     int   `json:"limit"`
}

func NewPaginationParams(page, limit int) *PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return &PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *PaginationParams) GetMeta(total int64) PaginationMeta {
	pageCount := int(total) / p.Limit
	if int(total)%p.Limit > 0 {
		pageCount++
	}

	return PaginationMeta{
		Total:     total,
		Page:      p.Page,
		PageCount: pageCount,
		Limit:     p.Limit,
	}
}
