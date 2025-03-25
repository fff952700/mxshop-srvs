package tests

import (
	"context"
	"google.golang.org/protobuf/encoding/protojson"
	"mxshop_srvs/order_srv/proto"
	"testing"
)

func TestCartCreate(t *testing.T) {
	resp, err := Client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		GoodsId: 2,
		Nums:    1,
		Checked: true,
	})
	if err != nil {
		t.Fatalf("CartItemCreate err: %v", err)
	}
	t.Logf("CartItemCreate resp: %v", resp)
}

func TestCartItemList(t *testing.T) {
	resp, err := Client.CartItemList(context.Background(), &proto.UserInfo{
		UserId: 1,
	})
	if err != nil {
		t.Fatalf("CartItemList err: %v", err)
	}
	// 使用proto序列化 否则会默认输出[]string格式
	data, err := protojson.Marshal(resp)
	if err != nil {
		t.Fatalf("CartItemList err: %v", err)
	}
	t.Logf("CartItemList resp: %v", string(data))
}

func TestCartItemUpdate(t *testing.T) {
	resp, err := Client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      3,
		Checked: false,
	})
	if err != nil {
		t.Fatalf("CartItemUpdate err: %v", err)
	}
	t.Logf("CartItemUpdate resp: %v", resp)
}

func TestDeleteCartItem(t *testing.T) {
	resp, err := Client.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		Id: 3,
	})
	if err != nil {
		t.Fatalf("CartItemDelete err: %v", err)
	}
	t.Logf("CartItemDelete resp: %v", resp)
}

func TestCreateOrder(t *testing.T) {
	resp, err := Client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京市海淀区",
		Name:    "张三",
		Mobile:  "13888888888",
		Post:    "test",
	})
	if err != nil {
		t.Fatalf("CreateOrder err: %v", err)
	}
	t.Logf("CreateOrder resp: %v", resp)
}
