package data

import (
	"strings"

	"github.com/looksaw/greenlight_2/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidatorFilter(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "page must greater than zero")
	v.Check(f.Page < 10000, "page", "page must less than 10000")
	v.Check(f.PageSize > 0, "page_size", "page size must greater than zero")
	v.Check(f.Page < 1000, "page_size", "page_size should less than 1000")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if safeValue == f.Sort {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsupported the sort field")
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
