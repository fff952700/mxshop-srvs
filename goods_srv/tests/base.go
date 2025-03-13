package tests

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/goods_srv/proto"
	"os"
	"time"
)

var GoodsClient proto.GoodsClient
var DB *gorm.DB

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	conn, err := grpc.NewClient("192.168.2.150:8091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panicw("err conn to server")
	}

	GoodsClient = proto.NewGoodsClient(conn)
	DbConn()
}

func DbConn() {
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
			//TablePrefix:   "mxshop_user_",
			NameReplacer: nil, //名称替换器（此处未使用）
		},
	})
	DB = db
	if err != nil {
		panic(err)
	}
}
