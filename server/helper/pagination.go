package helper

import "gorm.io/gorm"

func Pagination(limit, offset int) func(d *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Offset(offset).Limit(limit)
	}
}
