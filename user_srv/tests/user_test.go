package tests

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"mxshop_srvs/user_srv/proto"
	"strings"
	"testing"
)

func TestGetUserList(t *testing.T) {
	rsp, err := UserClient.GetUserList(context.Background(), &proto.PageInfo{
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
	// GORM 的 Omit 或 Select 配置、零值处理机制导致的。GORM 默认会忽略值为零的字段，而用数据库的默认值来填充。
	UserClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   "13888888888",
		Password: "123456",
	})
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
	fmt.Println(len(pwdInfo))
	// 解密
	verify := password.Verify("123456", pwdInfo[2], pwdInfo[3], options)
	fmt.Printf("verify:%v\n", verify)
	//
	result, err := UserClient.CheckUserPasswd(context.Background(), &proto.PasswordCheckInfo{
		Password:          "1234",
		EncryptedPassword: pwd,
	})
	if err != nil {
		panic(err)
	}
	//
	fmt.Printf("result:%v\n", result)

}
