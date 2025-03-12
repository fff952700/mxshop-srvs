package model

import "time"

type ShoppingCartItem struct {
	BaseModel
	GoodsId    int32   `gorm:"type:int;index:idx_goods_id;not null;comment:'商品id'"`
	GoodsImage string  `gorm:"type:varchar(200);not null;comment:'商品图片'"`
	GoodsName  string  `gorm:"type:varchar(200);not null;comment:'商品名称'"`
	GoodsPrice float32 `gorm:"type:decimal(10,2);not null;comment:'商品价格'"`
	Nums       int32   `gorm:"type:int;not null;default:1;comment:'购买数量'"`
	Checked    bool    `gorm:"type:tinyint(1);not null;default:1;comment:'是否勾选 1=已勾选,0=未勾选'"`
}

type OrderInfo struct {
	BaseModel
	UserId       int32     `gorm:"type:int;index;not null;comment:'用户id'"`
	OrderSn      string    `gorm:"type:varchar(30);not null;index;comment:'订单号'"`
	PayType      string    `gorm:"type:varchar(20);not null;default:alipay;comment:'支付类型，alipay:支付宝，wechat:微信'"`
	Status       string    `gorm:"type:varchar(20);comment:'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo      string    `gorm:"type:varchar(100);comment:'交易号'"`
	OrderMount   float32   `gorm:"type:decimal(10,2);comment:'订单金额'"`
	PayTime      time.Time `gorm:"type:datetime;comment:'支付时间'"`
	Address      string    `gorm:"type:varchar(150);comment:'收货地址'"`
	SignerName   string    `gorm:"type:varchar(20);comment:'签收人姓名'"`
	SingerMobile string    `gorm:"type:varchar(11);comment:'签收人电话'"`
	Post         string    `gorm:"type:varchar(20);comment"`
}
