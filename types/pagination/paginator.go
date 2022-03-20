package types

import (
	"math"
)

type Paginator struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	TotalRows  int    `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Table      string `json:"table"`
	QuerySet   []int  `json:"query_set"`
	// QuerySet   []db.VoicelineEntriesModel `json:"query_set"`QuerySet   []db.VoicelineEntriesModel `json:"query_set"`
}

func (p *Paginator) HasNext() bool {
	return (p.Page < p.NumberOfPages())
}

func (p *Paginator) NumberOfRows() (res int) {
	res = len(p.QuerySet)
	return
}

func (p *Paginator) NumberOfPages() int {
	return int(math.Ceil(float64(p.NumberOfRows()) / float64(p.Limit)))
}

func (p *Paginator) EndIndex() (res int) {
	if p.Page == p.NumberOfPages() {
		res = len(p.QuerySet)
		return
	}
	res = p.Page * p.Limit
	return
}

func (p *Paginator) NextPageNumber() int {
	return p.ValidateNumber(p.Page + 1)
}

func (p *Paginator) ValidateNumber(num int) int {
	if num < 1 || num == 1 {
		return 1
	}
	if num > p.NumberOfPages() {
		return 0
	}
	return num
}

func (p *Paginator) PreviousPageNumber() int {
	return p.ValidateNumber(p.Page - 1)
}

func (p *Paginator) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Paginator) Paginate() (res []int) {
	if len(p.QuerySet) <= p.Limit {
		return p.QuerySet
	}
	if p.Page == 1 {
		return p.QuerySet[0:p.Limit]
	}
	if !p.HasNext() {
		return p.QuerySet[p.PreviousPageNumber()*p.Limit:]
	}
	return p.QuerySet[p.EndIndex()-p.Limit : p.EndIndex()]
}
