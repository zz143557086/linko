package model

const (
	LEAVING_MESSAGES = iota + 1 // 留言类型: 1(留言)
	COMPLAINT                   // 留言类型: 2(投诉)
	INQUIRY                     // 留言类型: 3(询问)
	POST_SALE                   // 留言类型: 4(售后)
	WANT_TO_BUY                 // 留言类型: 5(求购)
)

type LeavingMessages struct {
	BaseModel // 留言信息基础模型

	User        int32  `gorm:"type:int;index;comment:'用户ID'"`                          // 用户ID
	MessageType int32  `gorm:"type:int;comment:'留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"` // 留言类型
	Subject     string `gorm:"type:varchar(100);comment:'主题'"`                         // 主题

	Message string `gorm:"type:text;comment:'留言内容'"`         // 留言内容
	File    string `gorm:"type:varchar(200);comment:'文件路径'"` // 文件路径
}

func (LeavingMessages) TableName() string {
	return "leavingmessages" // 数据库表名
}

type Address struct {
	BaseModel // 地址基础模型

	User         int32  `gorm:"type:int;index;comment:'用户ID'"`     // 用户ID
	Province     string `gorm:"type:varchar(10);comment:'省份'"`     // 省份
	City         string `gorm:"type:varchar(10);comment:'城市'"`     // 城市
	District     string `gorm:"type:varchar(20);comment:'区/县'"`    // 区/县
	Address      string `gorm:"type:varchar(100);comment:'详细地址'"`  // 详细地址
	SignerName   string `gorm:"type:varchar(20);comment:'收件人姓名'"`  // 收件人姓名
	SignerMobile string `gorm:"type:varchar(11);comment:'收件人手机号'"` // 收件人手机号
}

type UserFav struct {
	BaseModel // 用户收藏基础模型

	User  int32 `gorm:"type:int;index:idx_user_goods,unique;comment:'用户ID'"` // 用户ID
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique;comment:'商品ID'"` // 商品ID
}

func (UserFav) TableName() string {
	return "userfav" // 数据库表名
}
