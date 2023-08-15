package initialize

import (
	goredislib "github.com/go-redis/redis/v8"
	"shop_srvs/inventory_srv/global"
)

func InitRedis() {
	// 创建一个 Redis 客户端
	client := goredislib.NewClient(&goredislib.Options{
		Addr: global.Addr, // Redis 服务器的地址和端口号
	})

	global.RedisClinet = client
}
