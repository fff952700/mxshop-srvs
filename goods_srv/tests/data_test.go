package tests

import (
	"fmt"
	"math/rand"
	"mxshop_srvs/goods_srv/model"
	"testing"
)

func TestInsertBrand(t *testing.T) {
	var (
		brandList []*model.Brands
	)
	for i := 1; i < 30; i++ {
		brandList = append(brandList, &model.Brands{
			Name: fmt.Sprintf("brand_name_%d", i),
			Logo: fmt.Sprintf("https://brand.baidu.com/logo_%d.webp", i),
		})
	}
	DB.Create(&brandList)
}

func TestInsertCategory(t *testing.T) {
	var (
		topCategory   []*model.Category
		twoCategory   []*model.Category
		threeCategory []*model.Category
	)
	// 1级
	for i := 1; i <= 10; i++ {
		topCategory = append(topCategory, &model.Category{
			Name:             fmt.Sprintf("top_category_%d", i),
			ParentCategoryId: 0,
			Level:            1,
			IsTab:            true,
			Url:              fmt.Sprintf("https://search.baidu.com/category/top_category=%d", i),
		})
	}
	DB.Create(&topCategory)
	// 获取返回id
	for _, v := range topCategory {
		twoCategory = append(twoCategory, &model.Category{
			Name:             fmt.Sprintf("two_category_%d", v.ParentCategoryId),
			ParentCategoryId: v.Id,
			Level:            2,
			IsTab:            true,
			Url:              fmt.Sprintf("https://search.baidu.com/category/top_category=%d", v.Id),
		})
	}
	DB.Create(&twoCategory)

	for _, v := range twoCategory {
		threeCategory = append(threeCategory, &model.Category{
			Name:             fmt.Sprintf("three_category_%d", v.ParentCategoryId),
			ParentCategoryId: v.Id,
			Level:            3,
			IsTab:            true,
			Url:              fmt.Sprintf("https://search.baidu.com/category/three_category=%d", v.Id),
		})
	}
	DB.Create(&threeCategory)

}

func TestInsertGoods(t *testing.T) {
	var (
		goodsList []model.Goods
		category  model.Category
		brand     model.Brands
	)
	for i := 1; i < 30; i++ {
		DB.Order("RAND()").Limit(1).First(&category)
		DB.Order("RAND()").Limit(1).First(&brand)
		randInt := rand.Int31n(3) + 1
		goodsList = append(goodsList, model.Goods{
			CategoryId: category.Id,
			BrandId:    brand.Id,
			OnSale:
		})
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
