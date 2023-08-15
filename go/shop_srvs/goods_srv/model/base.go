package model

import (
	"database/sql/driver" // 引入数据库驱动程序包
	"encoding/json"       // 引入 JSON 相关的包
	"gorm.io/gorm"        // 引入 GORM ORM 框架包
	"time"                // 引入时间相关的包
)

type GormList []string

// 实现 driver.Valuer 接口，Value 方法将 GormList 转换为数据库支持的值
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g) // 将 GormList 转换为 JSON 字符串
}

// 实现 sql.Scanner 接口，Scan 方法将数据库存储的值转换为 GormList
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g) // 将数据库存储的 JSON 字符串转换为 GormList
}

type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int;comment:'id'" json:"id"`         // 记录ID字段，使用int32类型表示，指定为主键，并指定数据库字段类型为int
	CreatedAt time.Time      `gorm:"column:add_time;comment:'字段创建时间'" json:"-"`    // 记录创建时间字段，指定数据库列名为add_time
	UpdatedAt time.Time      `gorm:"column:update_time;comment:'字段更新时间'" json:"-"` // 记录更新时间字段，指定数据库列名为update_time
	DeletedAt gorm.DeletedAt `gorm:"comment:'字段删除时间'" json:"-"`                    // 记录删除时间字段，使用 GORM 框架提供的 gorm.DeletedAt 类型实现软删除功能
	IsDeleted bool           `gorm:"comment:'是否已经删除'" json:"-"`                    // 是否已经删除的标志字段
}
