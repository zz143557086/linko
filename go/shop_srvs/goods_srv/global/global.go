package global

import (
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

var (
	EsClient *elastic.Client
	DB       *gorm.DB

	Tags = []string{"linko", "lin", "goods_srv"}
)

const (
	Host             = "192.168.2.106" //数据库的ip
	Port             = "3306"          //数据库的端口
	User             = "root"          //数据库的账号密码
	Password         = "root"
	ServerPort       = 8004            //启动服务的ip
	ServerHost       = "127.0.0.1"     //端口
	ServerConsulHost = "192.168.2.106" //注册中心的ip
	SeverConsulPort  = 8500            //端口
	ServerName       = "goods_srv"     //启动服务的名称
)
