package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
	"strconv"
	"time"
)

// 生成MD5哈希值
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	// MySQL数据库连接信息
	dsn := "root:root@tcp(192.168.2.106:3306)/shop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	// 创建自定义logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 创建MySQL数据库连接实例
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// 自动创建数据库表
	_ = db.AutoMigrate(&model.Category{},
		&model.Brands{}, &model.GoodsCategoryBrand{}, &model.Banner{}, &model.Goods{})
	//Mysql2Es()
}

// 将MySQL数据导入Elasticsearch
func Mysql2Es() {
	// MySQL数据库连接信息
	dsn := "root:root@tcp(192.168.0.104:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	// 创建自定义logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 创建MySQL数据库连接实例
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// Elasticsearch连接信息
	host := "http://192.168.0.104:9200"
	logger := log.New(os.Stdout, "mxshop", log.LstdFlags)

	// 创建Elasticsearch客户端
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	var goods []model.Goods
	// 查询MySQL中的数据
	db.Find(&goods)
	for _, g := range goods {
		// 将MySQL的数据映射为Elasticsearch的数据模型
		esModel := model.EsGoods{
			ID:          g.ID,
			CategoryID:  g.CategoryID,
			BrandsID:    g.BrandsID,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			GoodsBrief:  g.GoodsBrief,
			ShopPrice:   g.ShopPrice,
		}

		// 将数据导入到Elasticsearch
		_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
		// 注意：在运行过程中，确保为Docker启动的Elasticsearch设置了足够的内存，否则可能会出
	}
}
