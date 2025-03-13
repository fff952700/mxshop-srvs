package tests

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_srvs/user_srv/proto"
)

var UserClient proto.UserClient

func init() {
	conn, err := grpc.NewClient("192.168.2.150:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	UserClient = proto.NewUserClient(conn)

}
