package handler

import "gorm.io/gorm"

// Paginate 分页函数
// page: 当前页码
// pageSize: 每页显示的数据条数
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 默认页码为1，若传入的页码为0，则默认为第一页
		if page == 0 {
			page = 1
		}

		// 设置每页显示的数据条数的范围
		switch {
		case pageSize > 100:
			// 如果传入的pageSize超过100，则将其设为100
			pageSize = 100
		case pageSize <= 0:
			// 如果传入的pageSize小于等于0，则将其设为10
			pageSize = 10
		}

		// 计算偏移量(offset)，即需要跳过的数据条数
		offset := (page - 1) * pageSize

		// 设置偏移量和每页数据条数，并返回结果
		return db.Offset(offset).Limit(pageSize)
	}
}
