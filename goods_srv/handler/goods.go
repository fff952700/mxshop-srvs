package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
	"reflect"
)

func (g *GoodsServer) Model2InfoResponse(goodsList interface{}) interface{} {
	switch goods := goodsList.(type) {
	case []*model.Goods:
		// 处理商品列表
		var goodsInfoResponse []*proto.GoodsInfoResponse
		for _, item := range goods {
			// 使用正确的函数调用，并传递 item（单个商品）
			goodsInfoResponse = append(goodsInfoResponse, g.toGoodsInfoResponse(item))
		}
		return goodsInfoResponse
	case *model.Goods:
		// 处理单个商品
		return g.toGoodsInfoResponse(goods)
	default:
		return nil
	}
}

// 提取公共转换函数
func (g *GoodsServer) toGoodsInfoResponse(goods *model.Goods) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              goods.Id,
		CategoryId:      goods.CategoryId,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		Images:          goods.Images,
		DescImages:      goods.DescImages,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
	}
}

func (g *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var (
		goodsList         []*model.Goods
		goodsListResponse = &proto.GoodsListResponse{}
	)

	// category query
	categoryQuery := global.DB.Model(&model.Category{})
	if req.TopCategory > 0 {
		categoryQuery = categoryQuery.Where("level = ?", req.TopCategory)
	}
	categoryQuery = categoryQuery.Where("is_tab = ?", req.IsTab)

	localDB := global.DB.Model(&model.Goods{})
	// 处理goods条件
	if req.PriceMax > 0 {
		localDB = localDB.Where("market_price <= ?", req.PriceMax)
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("market_price >= ?", req.PriceMin)
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brand_id = ?", req.Brand)
	}
	localDB = localDB.Where("is_new = ? AND is_hot = ?", req.IsNew, req.IsHot)

	// 子查询：使用 `categoryQuery` 作为子查询条件
	// 这里假设 categoryQuery 是从 `Category` 表中查询某些特定类别的商品的 id
	localDB = localDB.Where("category_id IN (?)", categoryQuery.Select("id"))

	// 添加分页
	localDB = localDB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums)))

	// 执行查询，将结果赋值给 GoodsList
	localDB.Find(&goodsList)
	if len(goodsList) < 1 {
		return nil, status.Error(codes.NotFound, "")
	}
	// 将查询结果转为返回格式
	goodsInfoResponse := g.Model2InfoResponse(goodsList)
	// 设置返回的响应数据
	goodsListResponse.Data = goodsInfoResponse.([]*proto.GoodsInfoResponse)
	goodsListResponse.Total = int32(len(goodsInfoResponse.([]*proto.GoodsInfoResponse)))
	return goodsListResponse, nil
}

// // 现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (g *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var (
		goodsList         []*model.Goods
		goodsListResponse = &proto.GoodsListResponse{}
	)
	if result := global.DB.Find(&goodsList, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "goods not found")
	}
	goodsListResponse.Total = int32(len(goodsList))
	goodsInfoResponses := g.Model2InfoResponse(goodsList)
	goodsListResponse.Data = goodsInfoResponses.([]*proto.GoodsInfoResponse)
	return goodsListResponse, nil
}

func (g *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var (
		goodsItem = &model.Goods{
			CategoryId:      req.CategoryId,
			BrandId:         req.BrandId,
			OnSale:          req.OnSale,
			ShipFree:        req.ShipFree,
			IsNew:           req.IsNew,
			IsHot:           req.IsHot,
			Name:            req.Name,
			GoodsSn:         req.GoodsSn,
			MarketPrice:     req.MarketPrice,
			ShopPrice:       req.ShopPrice,
			GoodsBrief:      req.GoodsBrief,
			GoodsFrontImage: req.GoodsFrontImage,
			Images:          req.Images,
			DescImages:      req.DescImages,
		}
	)

	if result := global.DB.Where("category_id = ? and brand_id = ? and name = ?", req.CategoryId, req.BrandId, req.Name).First(&goodsItem); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, " goods is exists")
	}
	global.DB.Create(goodsItem)
	goodsInfoResponse := g.toGoodsInfoResponse(goodsItem)
	return goodsInfoResponse, nil
}

func (g *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	var (
		goodsItem = &model.Goods{}
	)
	if result := global.DB.Where("id = ?", req.Id).Delete(&goodsItem); result.RowsAffected != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "delete goods err")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {

	var goodsItem = &model.Goods{}
	if result := global.DB.Where("id = ?", req.Id).First(&goodsItem); result.RowsAffected != 1 {
		return nil, status.Errorf(codes.NotFound, "goods not found")
	}

	v := reflect.ValueOf(req).Elem() // 获取 req 的指针值
	t := v.Type()

	// 用来记录需要更新的字段
	updates := map[string]interface{}{}
	// 遍历所有字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		// 跳过 id 字段，不进行更新
		if fieldName == "Id" {
			continue
		}

		// 判断字段是否为空，并准备更新
		if !field.IsZero() { // 如果字段不为空，则准备更新
			// 获取 goodsItem 的字段
			if goodsField := reflect.ValueOf(goodsItem).Elem().FieldByName(fieldName); goodsField.IsValid() && goodsField.CanSet() {
				// 将字段值加入更新列表
				updates[fieldName] = field.Interface()
			}
		}
	}

	// 如果有需要更新的字段，执行更新
	if len(updates) > 0 {
		if err := global.DB.Model(&goodsItem).Updates(updates).Error; err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update goods: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var (
		goodsItem = &model.Goods{}
	)
	if result := global.DB.Where("id = ?", req.Id).First(&goodsItem); result.RowsAffected != 1 {
		return nil, status.Errorf(codes.NotFound, " goods not found")
	}
	goodsInfoResponse := g.Model2InfoResponse(goodsItem)
	return goodsInfoResponse.(*proto.GoodsInfoResponse), nil

}
