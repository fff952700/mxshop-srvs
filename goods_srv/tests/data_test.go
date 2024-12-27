package tests

import (
	"fmt"
	"mxshop_srvs/goods_srv/model"
	"testing"
)

func TestInsertBrand(t *testing.T) {
	var (
		brandList []*model.Brands
		newBrands []*model.Brands
	)
	DB.Select("name,logo").Find(&brandList)
	for _, brands := range brandList {
		newBrands = append(newBrands, &model.Brands{
			Name: brands.Name,
			Logo: brands.Logo,
		})
	}
	DB.Table("new_brands").Create(&newBrands)

	//DB.Create(&brandList)
}

func TestInsertCategory(t *testing.T) {
	var (
		categoryList    []*model.Category
		newCategoryList []*model.Category
		//parentMap       = map[int32]int32{
		//	130364,
		//	130365
		//	130370
		//	135486
		//	135487
		//	135488
		//	135489
		//	136604
		//	136614
		//	136624
		//	136634
		//	136643
		//	136661
		//	136669
		//	136678
		//	136688
		//}
		resultList []struct {
			CategoryID    int `gorm:"column:category_id"`
			NewCategoryID int `gorm:"column:new_category_id"`
		}
	)
	DB.Raw("SELECT \n    c.id AS category_id,\n    c.name AS category_name,\n    n.id AS new_category_id,\n    n.name AS new_category_name\nFROM \n    category c\nLEFT JOIN \n    new_category n \nON \n    c.name = n.name\nWHERE \n    c.id IN (\n        SELECT \n            parent_category_id \n        FROM \n            category \n        WHERE \n            `level` = 3 \n        GROUP BY \n            parent_category_id\n    );").Scan(&resultList)
	parentMap := make(map[int32]int32)
	for _, result := range resultList {
		parentMap[int32(result.CategoryID)] = int32(result.NewCategoryID)
	}

	DB.Where("level = ?", 3).Find(&categoryList)
	//for _, category := range categoryList {
	//	fmt.Println(category.ParentCategoryId)
	//}
	//for _, category := range categoryList {
	//	newCategoryList = append(newCategoryList, &model.Category{
	//		Name:             category.Name,
	//		ParentCategoryId: category.ParentCategoryId,
	//		Level:            category.Level,
	//		IsTab:            category.IsTab,
	//	})
	//}
	for _, category := range categoryList {
		if newParentId, exists := parentMap[category.ParentCategoryId]; exists {
			//fmt.Printf("category %d, newParentId %d\n", category.ParentCategoryId, newParentId)
			newCategoryList = append(newCategoryList, &model.Category{
				Name:             category.Name,
				Level:            category.Level,
				IsTab:            category.IsTab,
				Url:              category.Url,
				ParentCategoryId: newParentId,
			})
		}
	}
	DB.Table("new_category").Create(&newCategoryList)

}

func TestInsertGoods(t *testing.T) {
	var (
		goodsList        []*model.Goods
		newGoodsList     []*model.Goods
		newCategoryBrand []*model.GoodsCategoryBrand
		categoryResult   []struct {
			CategoryID    int32 `gorm:"column:category_id"`
			NewCategoryID int32 `gorm:"column:new_category_id"`
		}
		brandResult []struct {
			BrandID    int32 `gorm:"column:brand_id"`
			NewBrandID int32 `gorm:"column:new_brand_id"`
		}
	)
	DB.Raw("SELECT \n    c.id AS category_id,\n    n.id AS new_category_id\nFROM \n    category c\nLEFT JOIN \n    new_category n \nON \n    c.name = n.name").Scan(&categoryResult)
	DB.Raw("SELECT \n    b.id AS brand_id,\n    n.id AS new_brand_id\nFROM \n    brands b\nLEFT JOIN \n    new_brands n \nON \n    b.name = n.name").Scan(&brandResult)
	categoryMap := make(map[int32]int32)
	for _, result := range categoryResult {
		categoryMap[result.CategoryID] = result.NewCategoryID
	}
	brandMap := make(map[int32]int32)
	for _, result := range brandResult {
		brandMap[result.BrandID] = result.NewBrandID
	}
	DB.Find(&goodsList)

	uniqueCategoryBrand := make(map[string]struct{})

	for _, g := range goodsList {
		// 跳过未找到分类或品牌的记录

		newCategoryId, categoryExists := categoryMap[g.CategoryId]
		newBrandId, brandExists := brandMap[g.BrandId]
		if !categoryExists || !brandExists {
			continue
		}
		// 构造唯一键用于去重
		uniqueKey := fmt.Sprintf("%d-%d", g.CategoryId, g.BrandId)
		if _, exists := uniqueCategoryBrand[uniqueKey]; !exists {
			uniqueCategoryBrand[uniqueKey] = struct{}{}
			newCategoryBrand = append(newCategoryBrand, &model.GoodsCategoryBrand{
				CategoryId: newCategoryId,
				BrandId:    newBrandId,
			})
		}

		// 构造新的 Goods 数据
		newGoodsList = append(newGoodsList, &model.Goods{
			CategoryId:      newCategoryId,
			BrandId:         newBrandId,
			OnSale:          g.OnSale,
			Name:            g.Name,
			ClickNum:        g.ClickNum,
			SoldNum:         g.SoldNum,
			FavNum:          g.FavNum,
			MarketPrice:     g.MarketPrice,
			ShopPrice:       g.ShopPrice,
			GoodsBrief:      g.GoodsBrief,
			ShipFree:        g.ShipFree,
			Images:          g.Images,
			DescImages:      g.DescImages,
			GoodsFrontImage: g.GoodsFrontImage,
			IsNew:           g.IsNew,
			IsHot:           g.IsHot,
		})
	}
	fmt.Println(len(newCategoryBrand))
	// 批量插入新数据
	if len(newGoodsList) > 0 {
		//DB.Table("new_goods").Create(&newGoodsList)
		DB.Table("new_goods_category_brand").Create(&newCategoryBrand)
	}

}

func TestInsertBanner(t *testing.T) {
	var (
		bannerList []*model.Banner
	)
	for i := 0; i < 3; i++ {
		bannerList = append(bannerList, &model.Banner{
			Image: fmt.Sprintf("https://image.baidu.com/%d", i),
			Url:   fmt.Sprintf("https://redirect.baidu.com/%d", i),
			Index: int32(i),
		})
	}
	DB.Create(&bannerList)

}
