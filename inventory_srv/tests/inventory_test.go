package tests

import (
	"context"
	"sync"
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
		GoodsId: 428,
		Stocks:  47,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestInvDetail(t *testing.T) {
	data, err := InventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 428,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestSell(t *testing.T) {
	// 测试并发访问 TODO 没有锁并发访问下获取的stocks可能相同
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			_, err := InventoryClient.Sell(context.Background(), &proto.SellInfo{
				GoodsInfo: []*proto.GoodsInvInfo{
					{GoodsId: 428, Stocks: 1},
				},
			})
			defer wg.Done()
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	t.Log("success")
}
