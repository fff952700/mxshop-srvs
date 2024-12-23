package model

// Category 商品分类表
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(255);" json:"name"`
	ParentCategoryId int32       `gorm:"type:int;DEFAULT 0;" json:"parent_category_id"` // 父级分类 ID
	SubCategory      []*Category `gorm:"-" json:"sub_category"`                         // 子分类（逻辑关联）
	Level            int32       `gorm:"type:tinyint; not null;comment '分类等级'" json:"level"`
	IsTab            bool        `gorm:"type:tinyint; not null;comment '是否显示'" json:"is_tab"`
	Url              string      `gorm:"type:varchar(255);" json:"url"`
}

// Brands 品牌表
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(255);not null"`
	Logo string `gorm:"type:varchar(255);not null;default ''"`
}

// GoodsCategoryBrand 商品分类与品牌的关联表
type GoodsCategoryBrand struct {
	BaseModel
	CategoryId int32 `gorm:"type:int; not null;uniqueIndex:uniq_category_brand"`
	BrandId    int32 `gorm:"type:int; not null;uniqueIndex:uniq_category_brand"`
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
	CategoryId      int32    `gorm:"type:int; not null;index:idx_category_id"` // 分类 ID，增加索引
	BrandId         int32    `gorm:"type:int; not null;index:idx_brand_id"`    // 品牌 ID，增加索引
	OnSale          bool     `gorm:"default false;not null;comment '是否上架'"`
	ShipFree        bool     `gorm:"default false;not null;comment '是否免运'"`
	IsNew           bool     `gorm:"default false;not null;comment '是否新商品'"`
	IsHot           bool     `gorm:"default false;not null;comment '是否热门'"`
	Name            string   `gorm:"type:varchar(255);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment '内部编号'"`
	ClickNum        int32    `gorm:"type:int;default 0;not null;comment '点击数'"`
	SoldNum         int32    `gorm:"type:int;default 0;not null;comment '购买护士'"`
	FavNum          int32    `gorm:"type:int;default 0;not null;comment '收藏数'"`
	MarketPrice     float32  `gorm:"not null;comment '商品价格'"`
	ShopPrice       float32  `gorm:"not null;comment '实际价格'"`
	GoodsBrief      string   `gorm:"type:varchar(255);not null;comment '商品介绍'"`
	GoodsFrontImage string   `gorm:"type:varchar(255);not null;comment '商品缩略图'"`
	Images          GormList `gorm:"type:json;not null;comment '商品图片'"`
	DescImages      GormList `gorm:"type:json;not null;comment '商品详情图片'"`
}
