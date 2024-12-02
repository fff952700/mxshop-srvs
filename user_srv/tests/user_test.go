package tests

import (
	"context"
	"crypto/sha512"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
)

var (
	userClient proto.UserClient
	db         *gorm.DB
)

func init() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)

}
func InitMysql() {
	dsn := "root:123456@tcp(localhost:23306)/mxshop?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql 阀值
			LogLevel:      logger.Info,
			Colorful:      true, //禁用色彩打印
		},
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //保持原有表名，不使用复数形式
			TablePrefix:   "mxshop_user_",
			NameReplacer:  nil, //名称替换器（此处未使用）
		},
	})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.User{})

}

func TestGetUserList(t *testing.T) {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil {
		panic(err)
	}
	if len(rsp.Data) < 1 {
		fmt.Println("没有用户")
		return
	}
	for _, value := range rsp.Data {
		fmt.Println(value)
	}
}

// 注册用户
func TestCreateUser(t *testing.T) {
	InitMysql()
	var user model.User
	user.Nickname = "test"
	user.Mobile = "13888888888"
	user.Role = 1
	// GORM 的 Omit 或 Select 配置、零值处理机制导致的。GORM 默认会忽略值为零的字段，而用数据库的默认值来填充。
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("123456", options)
	pwd := fmt.Sprintf("$sha512$%s$%s", salt, encodedPwd)
	user.Password = pwd
	fmt.Println(user)
	tx := db.Create(&user)
	if tx.Error != nil {
		fmt.Printf("创建用户失败 %v", tx.Error)
	}
}

// 校验密码
func TestPassWordCheck(t *testing.T) {
	InPwd := "123456"
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(InPwd, options)
	// 秘文
	pwd := fmt.Sprintf("$sha512$%s$%s", salt, encodedPwd)
	fmt.Println(pwd)
	pwdInfo := strings.Split(pwd, "$")
	// 解密
	verify := password.Verify("1111", pwdInfo[2], pwdInfo[3], options)
	fmt.Printf("verify:%v\n", verify)
	//
	result, err := userClient.CheckUserPasswd(context.Background(), &proto.PasswordCheckInfo{
		Password:          "1234",
		EncryptedPassword: pwd,
	})
	if err != nil {
		panic(err)
	}
	//
	fmt.Printf("result:%v\n", result)

}
