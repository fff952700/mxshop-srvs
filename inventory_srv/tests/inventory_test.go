package tests

import (
	"context"
	"testing"

	"math/rand"

	model2 "mxshop_srvs/goods_srv/model"
	"mxshop_srvs/inventory_srv/model"
	"mxshop_srvs/inventory_srv/proto"
)

func TestInitAddInventory(t *testing.T) {
	var (
		InventoryList []*model.Inventory
		GoodsList     []*model2.Goods
	)
	if result := e.goodsDB.Find(&GoodsList); result.RowsAffected == 0 {
		return
	}

	for _, v := range GoodsList {
		// 每次循环创建新的 Inventory 实例
		inventory := &model.Inventory{
			GoodsId: v.Id,
			Stocks:  rand.Int31n(100),
		}
		InventoryList = append(InventoryList, inventory)
	}
	e.inventoryDB.Save(&InventoryList)
}

func TestSetInv(t *testing.T) {
	_, err := InventoryClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 1,
		Stocks:  345,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestInvDetail(t *testing.T) {
	data, err := InventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 1,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}
