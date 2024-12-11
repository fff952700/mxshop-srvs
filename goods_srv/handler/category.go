package handler

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	var categoryList []model.Category
	// SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	// 配置指明了外键后，可以使用Preload预加载，来把品牌的子分类也取出来
	result := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categoryList)
	b, _ := json.Marshal(&categoryList)
	zap.S().Info(string(b))
	return &proto.CategoryListResponse{Total: int32(result.RowsAffected), JsonData: string(b)}, nil
}

func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	// 判断分类是否存在
	var category model.Category
	if result := global.DB.Where(model.Category{Level: req.Level}).First(&category); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}
	// 实例化返回对象处理1级分类
	categoryListResponse := &proto.SubCategoryListResponse{}
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.Id,
		Name:           category.Name,
		Level:          category.Level,
		ParentCategory: category.ParentCategoryId,
		IsTab:          category.IsTab,
	}
	// 1及分类对象
	zap.S().Infof("top category %v", category)
	// 处理子类
	var SubCategoryList []model.Category
	var SubCateGoryResponseList []*proto.CategoryInfoResponse
	Preload := "SubCategory"
	if req.Level == 1 {
		Preload = "SubCategory.SubCategory"
	}
	global.DB.Where(model.Category{Level: 1}).Preload(Preload).Find(&SubCategoryList)
	for _, SubCategory := range SubCategoryList {
		SubCateGoryResponseList = append(SubCateGoryResponseList, &proto.CategoryInfoResponse{
			Id:             SubCategory.Id,
			Name:           SubCategory.Name,
			Level:          SubCategory.Level,
			ParentCategory: SubCategory.ParentCategoryId,
			IsTab:          SubCategory.IsTab,
		})
	}
	zap.S().Infof("two category %v", SubCateGoryResponseList)
	categoryListResponse.SubCategory = SubCateGoryResponseList
	return categoryListResponse, nil
}
