package model

type Inventory struct {
	BaseModel
	GoodsId  string `gorm:"type:int;unique;NOT NULL;"`
	GoodsNum int    `gorm:"type:int;NOT NULL;"`
}

type InventoryHistory struct {
	BaseModel
	GoodsId  string `gorm:"type:int;NOT NULL;"`
	GoodsNum int    `gorm:"type:int;NOT NULL;"`
	OrderId  string `gorm:"type:int;NOT NULL;"`
	Status   string `gorm:"type:tinyint(1);comment '1 预扣 2 已完成 3 已取消';NOT NULL;"`
}
