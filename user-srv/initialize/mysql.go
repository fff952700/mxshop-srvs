package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop-srvs/user-srv/global"
	"os"
	"time"
)

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.ServerConf.MysqlInfo.Username, global.ServerConf.MysqlInfo.Password, global.ServerConf.MysqlInfo.Host, global.ServerConf.MysqlInfo.Port, global.ServerConf.MysqlInfo.Dbname)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql 阀值
			LogLevel:      logger.Info,
			Colorful:      true, //禁用色彩打印
		},
	)

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //保持原有表名，不使用复数形式
			TablePrefix:   "mxshop_user_",
			NameReplacer:  nil, //名称替换器（此处未使用）
		},
	})
	global.DB = DB
	if err != nil {
		panic(err)
	}
}
