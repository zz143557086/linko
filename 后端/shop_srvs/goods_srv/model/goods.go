package model

// 定义了商品分类的结构体
// Category 分类表
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null;comment:'分类名称'" json:"name"`                           // 分类名称，最大长度为20
	ParentCategoryID int32       `json:"parent" gorm:"comment:'父分类ID'" `                                                 // 父分类ID
	ParentCategory   *Category   `json:"-" gorm:"comment:'父分类'" `                                                        // 父分类
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID; comment:'子分类集合'" json:"sub_category"` // 子分类集合
	Level            int32       `gorm:"type:int;not null;default:1;comment:'分类级别'" json:"level"`                        // 分类级别
	IsTab            bool        `gorm:"default:false;not null;comment:'是否显示在导航栏标签'" json:"is_tab" `                     // 是否显示在导航栏标签
}

// Brands 品牌表
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null;comment:'品牌名称'" `                   // 品牌名称，最大长度为20
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:'品牌logo图片链接'" ` // 品牌logo图片链接
}

// GoodsCategoryBrand 商品分类和品牌关联表
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique;comment:'分类ID'" ` // 分类ID
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique;comment:'品牌ID'" ` // 品牌ID
	Brands   Brands
}

// 返回 GoodsCategoryBrand 表的表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// Banner 轮播图表
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null;comment:'轮播图图片链接'" `   // 轮播图图片链接
	Url   string `gorm:"type:varchar(200);not null;comment:'轮播图点击跳转链接'" ` // 轮播图点击跳转链接
	Index int32  `gorm:"type:int;default:1;not null;comment:'轮播图显示顺序'" `  // 轮播图显示顺序
}

// Goods 商品表
type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null;comment:'分类ID'" ` // 分类ID
	Category   Category

	BrandsID int32 `gorm:"type:int;not null;comment:'品牌ID'" ` // 品牌ID
	Brands   Brands

	OnSale          bool     `gorm:"default:false;not null;comment:'是否上架'" `         // 是否上架
	ShipFree        bool     `gorm:"default:false;not null;comment:'是否包邮'" `         // 是否包邮
	IsNew           bool     `gorm:"default:false;not null;comment:'是否新品'" `         // 是否新品
	IsHot           bool     `gorm:"default:false;not null;comment:'是否热销'" `         // 是否热销
	Name            string   `gorm:"type:varchar(50);not null;comment:'商品名称'" `      // 商品名称，最大长度为50
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:'商品编号'" `      // 商品编号，最大长度为50
	ClickNum        int32    `gorm:"type:int;default:0;not null;comment:'点击量'" `     // 点击量
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:'销量'"`       // 销量
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:'收藏量'" `     // 收藏量
	MarketPrice     float32  `gorm:"not null;comment:'市场价'"`                         // 市场价
	ShopPrice       float32  `gorm:"not null;comment:'销售价'"`                         // 销售价
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:'商品简介'"`      // 商品简介，最大长度为100
	Images          GormList `gorm:"type:varchar(1000);not null;comment:'商品图片列表'"`   // 商品图片列表
	DescImages      GormList `gorm:"type:varchar(1000);not null;comment:'商品详情图片列表'"` // 商品详情图片列表
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:'商品封面图片链接'"`  // 商品封面图片链接
}

/*
// 在创建记录之后执行的回调函数
func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	// 构造要保存到 Elasticsearch 中的模型
	esModel := EsGoods{
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

	// 将模型保存到 Elasticsearch 中
	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}


// 在删除记录之后执行的回调函数
func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	// 从 Elasticsearch 中删除对应的模型
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}*/
