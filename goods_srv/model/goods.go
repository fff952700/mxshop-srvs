package model

// Category 商品分类表
type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(255);"`
	ParentCategoryId int32  `gorm:"type:int; not null"`
	ParentCategory   *Category
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryId;references:Id"` // 通过 ParentCategoryId 关联到其他 Category 记录的 Id 字段
	Level            int32       `gorm:"type:tinyint; not null;comment '分类等级'"`
	IsTab            bool        `gorm:"type:tinyint; not null;comment '是否显示'"`
}

// Brands 品牌表
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(255);not null"`
	Logo string `gorm:"type:varchar(255);not null;default ''"`
}

// GoodsCategoryBrand modelToModel
type GoodsCategoryBrand struct {
	BaseModel
	CategoryId int32    `gorm:"type:int; not null;index:idx_category_brand,unique"`
	Category   Category `gorm:"foreignKey:CategoryId"`
	BrandId    int32    `gorm:"type:int; not null;index:idx_category_brand,unique"`
	Brands     Brands   `gorm:"foreignKey:BrandId"`
}

// Banner 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(255);not null;comment 'url地址'"`
	Url   string `gorm:"type:varchar(255);not null;comment 'redirect url'"`
	Index int32  `gorm:"type:int;not null;comment 'banner 顺序'"`
}

// Goods 商品表
type Goods struct {
	BaseModel
	CategoryId      int32    `gorm:"type:int; not null;"`
	Category        Category `gorm:"foreignKey:CategoryId"`
	BrandId         int32    `gorm:"type:int; not null;"`
	Brands          Brands   `gorm:"foreignKey:BrandId"`
	OnSale          bool     `gorm:"default false;not null;comment '是否上架'"`
	ShipFree        bool     `gorm:"default false;not null;comment '是否免运'"`
	IsNew           bool     `gorm:"default false;not null;comment '是否新商品'"`
	IsHot           bool     `gorm:"default false;not null;comment '是否热门'"`
	Name            string   `gorm:"type:varchar(255);not null"`
	GoodsSn         string   `gorm:"type:int;not null;comment '内部编号'"`
	CheckNum        int32    `gorm:"type:int;default 0;not null;comment '点击数'"`
	SoldNum         int32    `gorm:"type:int;default 0;not null;comment '购买护士'"`
	FavNum          int32    `gorm:"type:int;default 0;not null;comment '收藏数'"`
	MarketPrice     float32  `gorm:"not null;comment '商品价格'"`
	ShopPrice       float32  `gorm:"not null;comment '实际价格'"`
	GoodsBrief      string   `gorm:"type:varchar(255);not null;comment '商品介绍'"`
	GoodsFrontImage string   `gorm:"type:varchar(255);not null;comment '商品缩略图'"`
	Image           GormList `gorm:"type:varchar(1000);not null;comment '商品图片'"`
	DescImages      GormList `gorm:"type:varchar(1000);not null;comment '商品详情图片'"`
}
