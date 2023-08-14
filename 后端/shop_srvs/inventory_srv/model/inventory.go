package model

import (
	"database/sql/driver"
	"encoding/json"
)

//type Stock struct {
//	BaseModel
//	Name string
//	Address string
//}

type GoodsDetail struct {
	Goods int32 `gorm:"column:goods;comment:商品ID"` // 商品ID
	Num   int32 `gorm:"column:num;comment:数量"`     // 数量
}
type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"column:goods;type:int;index;comment:商品ID"` // 商品ID
	Stocks  int32 `gorm:"column:stocks;type:int;comment:库存"`        // 库存
	Version int32 `gorm:"column:version;type:int;comment:分布式锁的乐观锁"` // 分布式锁的乐观锁
}

type InventoryNew struct {
	BaseModel
	Goods   int32 `gorm:"column:goods;type:int;index;comment:商品ID"` // 商品ID
	Stocks  int32 `gorm:"column:stocks;type:int;comment:库存"`        // 库存
	Version int32 `gorm:"column:version;type:int;comment:分布式锁的乐观锁"` // 分布式锁的乐观锁
	Freeze  int32 `gorm:"column:freeze;type:int;comment:冻结库存"`      // 冻结库存
}

type Delivery struct {
	Goods   int32  `gorm:"column:goods;type:int;index;comment:商品ID"`                             // 商品ID
	Nums    int32  `gorm:"column:nums;type:int;comment:数量"`                                      // 数量
	OrderSn string `gorm:"column:order_sn;type:varchar(200);comment:订单编号"`                       // 订单编号
	Status  string `gorm:"column:status;type:varchar(200);comment:状态：1. 表示等待支付，2. 表示支付成功，3. 失败"` // 状态：1. 表示等待支付，2. 表示支付成功，3. 失败
}

type StockSellDetail struct {
	OrderSn string          `gorm:"column:order_sn;type:varchar(200);index:idx_order_sn,unique;comment:订单编号"` // 订单编号
	Status  int32           `gorm:"column:status;type:varchar(200);comment:状态；1 表示已扣减，2. 表示已归还"`              // 状态：1 表示已扣减，2. 表示已归还
	Detail  GoodsDetailList `gorm:"column:detail;type:varchar(200);comment:商品明细"`                             // 商品明细
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
} // 库存销售详情表

//type InventoryHistory struct {
//	user int32
//	goods int32
//	nums int32
//	order int32
//	status int32 //1. 表示库存是预扣减， 幂等性， 2. 表示已经支付
//}
