package main

import (
	"log"
	"mxshop_srvs/order_srv/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := "root:123456@tcp(localhost:23306)/mxshop_order?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql 阀值
			LogLevel:      logger.Info,
			Colorful:      true, //禁用色彩打印
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //保持原有表名，不使用复数形式
			//TablePrefix:   "goods_",
			NameReplacer: nil, //名称替换器（此处未使用）
		},
	})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.OrderGoods{}, &model.ShoppingCart{}, &model.OrderInfo{})
}
