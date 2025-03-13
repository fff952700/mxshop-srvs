package tests

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mxshop_srvs/inventory_srv/proto"
)

var (
	InventoryClient proto.InventoryClient
	e               *engine
)

type engine struct {
	goodsDB     *gorm.DB
	inventoryDB *gorm.DB
}

func newEngine() *engine {
	return &engine{
		goodsDB:     nil,
		inventoryDB: nil,
	}
}

func init() {
	e = newEngine()
	loggers, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(loggers)

	conn, err := grpc.NewClient("192.168.2.150:8092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panicw("err conn to server")
	}
	InventoryClient = proto.NewInventoryClient(conn)
	e.goodsDB = GoodsDB()
	e.inventoryDB = InventoryDB()
}

func GoodsDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"192.168.2.150",
		23306,
		"mxshop_goods")
	zap.S().Infof("dsn:%s", dsn)
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
			NameReplacer:  nil,  //名称替换器（此处未使用）
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	return db
}

func InventoryDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"192.168.2.150",
		23306,
		"mxshop_inventory")
	zap.S().Infof("dsn:%s", dsn)
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
			NameReplacer:  nil,  //名称替换器（此处未使用）
		},
	})
	if err != nil {
		panic(err)
	}
	return db
}
