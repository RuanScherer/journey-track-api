package scopes

import "gorm.io/gorm"

type PaginationOptions struct {
	Page     int
	PageSize int
}

func Paginate(options PaginationOptions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := options.Page
		if page <= 0 {
			page = 1
		}

		pageSize := options.PageSize
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
