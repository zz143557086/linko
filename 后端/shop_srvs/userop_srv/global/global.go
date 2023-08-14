package global

import (
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Tags = []string{"linko", "lin", "userop_srv"}
)

const (
	Host             = "192.168.2.106" //数据库
	Port             = "3306"
	User             = "root"
	Password         = "root"
	ServerPort       = 8007 //自己的服务
	ServerHost       = "127.0.0.1"
	ServerConsulHost = "192.168.2.106" //注册中心
	SeverConsulPort  = 8500
	ServerName       = "userop_srv"
)

//func init() {
//	dsn := "root:root@tcp(192.168.0.104:3306)/shop_userop_srv?charset=utf8mb4&parseTime=True&loc=Local"
//
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//		logger.Config{
//			SlowThreshold: time.Second,   // 慢 SQL 阈值
//			LogLevel:      logger.Info, // Log level
//			Colorful:      true,         // 禁用彩色打印
//		},
//	)
//
//	// 全局模式
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic(err)
//	}
//}
