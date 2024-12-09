package tests

import (
	"go.uber.org/zap"
	"mxshop_srvs/goods_srv/model"
	"testing"
)

func TestGetCateGoryAll(t *testing.T) {
	InitMysql()
	var CategoryList []model.Category
	// 只取1个SubCategory 只会加载到二级分类
	DB.Preload("SubCategory.SubCategory").Where(model.Category{Level: 1}).Find(&CategoryList)
	for _, category := range CategoryList {
		zap.S().Infof("Category List: %v", category)
	}
}
