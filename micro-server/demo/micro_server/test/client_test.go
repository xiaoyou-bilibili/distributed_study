package test

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	pb "server/api"
	"testing"
)

func TestClient(t *testing.T) {
	config := api.DefaultConfig()
	config.Address = "192.168.1.40:30571"
	consulClient, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	// 注册consul服务
	r := consul.New(consulClient)
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///kratos-server"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
		))
	fmt.Println(err)
	// 建立连接
	a := pb.NewUserClient(conn)
	res, err := a.GetInfo(context.Background(), &pb.UserRequest{Id: 2})
	fmt.Println(err)
	fmt.Println("结果", res)
}
