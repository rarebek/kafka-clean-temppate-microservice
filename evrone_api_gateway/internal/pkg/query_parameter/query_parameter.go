package query_parameter

import (
	"strconv"
)

type QueryParameter interface {
	GetLimit() uint64
	GetPage() uint64
	GetOffset() uint64
	GetParameters() map[string]string
}

func New(values map[string][]string) QueryParameter {
	qp := parameters{
		limit:  10,
		page:   1,
		offset: 0,
		values: make(map[string]string, 10),
	}

	for key, val := range values {
		if key == "limit" {
			if limit, err := strconv.ParseUint(val[0], 10, 64); err == nil {
				qp.limit = limit
			}
			continue
		}
		if key == "page" {
			if page, err := strconv.ParseUint(val[0], 10, 64); err == nil {
				qp.page = page
			}
			continue
		}
		if key == "offset" {
			if offset, err := strconv.ParseUint(val[0], 10, 64); err == nil {
				qp.offset = offset
			}
			continue
		}
		qp.values[key] = val[0]
	}

	return &qp
}

type parameters struct {
	limit  uint64
	page   uint64
	offset uint64
	values map[string]string
}

func (p *parameters) GetParameters() map[string]string {
	params := make(map[string]string)
	for key, val := range p.values {
		if len(val) >= 1 {
			params[key] = val
		}
	}
	return params
}

func (p *parameters) GetLimit() uint64 {
	return p.limit
}

func (p *parameters) GetPage() uint64 {
	return p.page
}

func (p *parameters) GetOffset() uint64 {
	return p.offset
}
