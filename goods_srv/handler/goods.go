package handler

import (
	"context"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

func (g *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var (
		GoodsList         []model.Goods
		GoodsInfoResponse []*proto.GoodsInfoResponse
		GoodsListResponse = &proto.GoodsListResponse{}
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
	localDB.Find(&GoodsList)

	// 将查询结果转为返回格式
	for _, goods := range GoodsList {
		GoodsInfoResponse = append(GoodsInfoResponse, &proto.GoodsInfoResponse{
			Id:              goods.Id,
			CategoryId:      goods.CategoryId,
			Name:            goods.Name,
			GoodsSn:         goods.GoodsSn,
			ClickNum:        goods.CheckNum,
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
		})
	}

	// 设置返回的响应数据
	GoodsListResponse.Data = GoodsInfoResponse
	GoodsListResponse.Total = int32(len(GoodsInfoResponse))
	return GoodsListResponse, nil
}

// // 现在用户提交订单有多个商品，你得批量查询商品的信息吧
//func (g *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method BatchGetGoods not implemented")
//}
//func (g *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method CreateGoods not implemented")
//}
//func (g *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method DeleteGoods not implemented")
//}
//func (g *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method UpdateGoods not implemented")
//}
//func (g *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method GetGoodsDetail not implemented")
//}
