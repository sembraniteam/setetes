package response

import (
	"math"
	"regexp"
	"strings"
)

type (
	Pagination struct {
		page   int    `form:"page" validate:"omitempty,numeric,min=1"`
		limit  int    `form:"limit" validate:"omitempty,numeric,min=1,max=1000"`
		search string `form:"search" validate:"omitempty,max=100"`
		sort   string `form:"sort" validate:"omitempty,oneof=ASC DESC" reason:"oneof=Order must be one of ASC, DESC"`
	}

	BaseEntries[T any] struct {
		Entries       []T   `json:"entries"`
		HasReachedMax bool  `json:"has_reached_max"`
		TotalPages    int64 `json:"total_pages"`
	}
)

func Entries[T any](entries []T, hasReachedMax bool, totalPages int64) *BaseEntries[T] {
	return &BaseEntries[T]{
		Entries:       entries,
		HasReachedMax: hasReachedMax,
		TotalPages:    totalPages,
	}
}

func (p *Pagination) GetPage() int {
	if p.page == 0 {
		p.page = 1
	}

	return p.page
}

func (p *Pagination) GetLimit() int {
	if p.limit == 0 {
		p.limit = 10
	}

	return p.limit
}

func (p *Pagination) GetSearch() string {
	s := strings.ToLower(strings.TrimSpace(p.search))
	if (len(s)) > 100 {
		s = s[:100]
	}

	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = regexp.MustCompile(`[^\w\s\-]`).ReplaceAllString(s, "")

	return s
}

func (p *Pagination) GetSort() string {
	if p.sort == "" {
		p.sort = "DESC"
	}

	return p.sort
}

func (p *Pagination) GetTotalPages(items int64) int64 {
	if items == 0 {
		return 1
	}

	return int64(math.Ceil(float64(items) / float64(p.GetLimit())))
}

func (p *Pagination) GetHasReachedMax(items int64) bool {
	return int64(p.GetPage()) >= p.GetTotalPages(items)
}
