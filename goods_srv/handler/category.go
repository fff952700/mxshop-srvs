package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

/*
GetAllCategoryList
获取所有分类，需要返回分类和子分类的所有信息 需要用到外键
GORM 可以通过 Preload 预加载 has many 关联的记录，查看 预加载 获取详情
*/
func (g *GoodsServer) GetAllCategoryList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var CategoryList []model.Category
	// 只取1个SubCategory 只会加载到二级分类
	global.DB.Preload("SubCategory.SubCategory").Where(model.Category{Level: 1}).Find(&CategoryList)
	for _, category := range CategoryList {
		zap.S().Infof("Category List: %v", category)
	}
	return nil, nil
}
