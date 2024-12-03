package tests

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/goods_srv/handler"
	"os"
	"testing"
	"time"

	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

var DB *gorm.DB

func TestGetBrandList(t *testing.T) {
	// 实例化查询对象
	InitMysql()
	var brands []model.Brands
	// 分页返回数据 result存在元数据 将数据赋值给brands是通过Find(&brands)完成的
	result := DB.Scopes(handler.Paginate(1, 20)).Find(&brands)
	fmt.Println(brands)
	// 获取总条数
	brandResponses := make([]*proto.BrandInfoResponse, 0)
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.Id,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	fmt.Println(result.RowsAffected, brandResponses)
}

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"192.168.2.106",
		23306,
		"mxshop_goods")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql 阀值
			LogLevel:      logger.Info,
			Colorful:      true, //禁用色彩打印
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //保持原有表名，不使用复数形式
			//TablePrefix:   "mxshop_user_",
			NameReplacer: nil, //名称替换器（此处未使用）
		},
	})

	if err != nil {
		panic(err)
	}
}
