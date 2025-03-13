package tests

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_srvs/order_srv/proto"
)

var (
	//e      *engine
	Client proto.OrderClient
)

//type engine struct {
//	mysql *gorm.DB
//}

func init() {
	//e = newEngine()
	loggers, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(loggers)

	conn, err := grpc.NewClient("192.168.2.150:8093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panicw("err conn to server")
	}
	Client = proto.NewOrderClient(conn)
}

//func newEngine() *engine {
//	return &engine{
//		mysql: nil,
//	}
//}
