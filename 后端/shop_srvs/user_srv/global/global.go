package global

import (
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Tags = []string{"linko", "lin", "user_srv"}
)

const (
	Host             = "192.168.2.106" //数据库
	Port             = "3306"
	User             = "root"
	Password         = "root"
	ServerPort       = 8003 //自己的服务
	ServerHost       = "127.0.0.1"
	ServerConsulHost = "192.168.2.106" //注册中心
	SeverConsulPort  = 8500
	ServerName       = "user_srv"
)
