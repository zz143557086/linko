package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

/*
1. 密文 2. 密文不可反解
 1. 对称加密
 2. 非对称加密
 3. md5 信息摘要算法
    密码如果不可以反解，用户找回密码
*/
type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11) comment '手机号码';not null"`
	Password string     `gorm:"type:varchar(100) comment '密码';not null"`
	Name     string     `gorm:"type:varchar(20) comment '姓名'"`
	Birthday *time.Time `gorm:"type:datetime comment '生日'"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role     int32      `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}
