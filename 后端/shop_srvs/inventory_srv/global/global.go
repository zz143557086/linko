package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	Tags        = []string{"linko", "lin", "inventory_srv"}
	RedisClinet *redis.Client // 分布式锁
)

const (
	Host             = "192.168.2.106" //数据库
	Port             = "3306"
	User             = "root"
	Password         = "root"
	ServerPort       = 8005 //自己启动的服务
	ServerHost       = "127.0.0.1"
	ServerConsulHost = "192.168.2.106" //注册中心
	SeverConsulPort  = 8500
	ServerName       = "inventory_srv"
	Addr             = "192.168.2.106:6379" //redis
	RocketMQ         = "192.168.2.106:9876"
)
