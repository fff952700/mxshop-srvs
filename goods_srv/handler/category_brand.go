package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

// CategoryBrandList 获取所有品牌分类
func (g *GoodsServer) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	// 去掉外键后需要手动构建返回数据
	var (
		categoryBrandList         []model.GoodsCategoryBrand
		categoryBrandResp         []*proto.CategoryBrandResponse
		categoryMap               = make(map[int32]model.Category)
		brandMap                  = make(map[int32]model.Brands)
		categoryBrandListResponse = &proto.CategoryBrandListResponse{}
		categories                []model.Category
		brands                    []model.Brands
		categoryIDs               []int32
		brandIDs                  []int32
	)
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&categoryBrandList)
	categoryBrandListResponse.Total = int32(result.RowsAffected)
	// 收集categoryId,brandId
	for _, cb := range categoryBrandList {
		categoryIDs = append(categoryIDs, cb.CategoryId)
		brandIDs = append(brandIDs, cb.BrandId)
	}

	//获取到查询结果
	global.DB.Where("id in ?", categoryIDs).Find(&categories)
	global.DB.Where("id in ?", brandIDs).Find(&brands)
	// 构造 map 以便快速查找
	for _, category := range categories {
		categoryMap[category.Id] = category
	}
	for _, brand := range brands {
		brandMap[brand.Id] = brand
	}
	// 构建返回对象
	for _, categoryBrand := range categoryBrandList {
		category, catOk := categoryMap[categoryBrand.CategoryId]
		brand, brandOk := brandMap[categoryBrand.BrandId]
		if catOk && brandOk {
			categoryBrandResp = append(categoryBrandResp, &proto.CategoryBrandResponse{
				Category: &proto.CategoryInfoResponse{
					Id:             category.Id,
					Name:           category.Name,
					ParentCategory: category.ParentCategoryId,
					Level:          category.Level,
					IsTab:          category.IsTab,
				},
				Brand: &proto.BrandInfoResponse{
					Id:   brand.Id,
					Name: brand.Name,
					Logo: brand.Logo,
				},
			})
		}
	}
	categoryBrandListResponse.Data = categoryBrandResp
	return categoryBrandListResponse, nil

}

// GetCategoryBrandList 通过category获取brands
func (s *GoodsServer) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	var (
		brandListResponse = &proto.BrandListResponse{}
		categoryBrandList []model.GoodsCategoryBrand
		brandIds          []int32
		brands            []model.Brands
		brandMap          = make(map[int32]model.Brands)
		brandInfoResponse []*proto.BrandInfoResponse
	)
	// 查看分类是否存在
	result := global.DB.Where("category_id = ?", req.Id).Find(&categoryBrandList)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}
	brandListResponse.Total = int32(result.RowsAffected)
	// 获取所有brandId
	for _, cb := range categoryBrandList {
		brandIds = append(brandIds, cb.BrandId)
	}
	global.DB.Where("id in ?", brandIds).Find(&brands)
	for _, brand := range brands {
		brandMap[brand.Id] = brand
	}
	// 构造返回数据
	for _, cb := range categoryBrandList {
		if brand, ok := brandMap[cb.BrandId]; ok {
			brandInfoResponse = append(brandInfoResponse, &proto.BrandInfoResponse{
				Id:   brand.Id,
				Name: brand.Name,
				Logo: brand.Logo,
			})
		}
	}
	brandListResponse.Data = brandInfoResponse
	return brandListResponse, nil
}

func (g *GoodsServer) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var (
		categoryBrand = model.GoodsCategoryBrand{
			CategoryId: req.CategoryId,
			BrandId:    req.BrandId,
		}
	)
	if result := global.DB.First(&model.Category{}, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}
	if result := global.DB.First(&model.Brands{}, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "brands not found")
	}
	global.DB.Create(&categoryBrand)
	return &proto.CategoryBrandResponse{Id: categoryBrand.Id}, nil
}

//func (g *GoodsServer) DeleteCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
//
//}
//func (g *GoodsServer) UpdateCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
//
//}
