package service

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"server/ent"
	"server/ent/user"

	pb "server/api"
)

var client *ent.Client

func init() {
	var err error
	// 建立数据库连接
	client, err = ent.Open("mysql", "root:xiaoyou@tcp(192.168.1.40:30000)/micro?parseTime=True")
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Close()
	// 自动创建表
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

type UserService struct {
	pb.UnimplementedUserServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetInfo(ctx context.Context, req *pb.UserRequest) (*pb.UserInfo, error) {
	// 查询数据
	u, err := client.User.
		Query().Where(user.ID(int(req.Id))).
		Only(context.Background())
	if err != nil {
		// 打印错误信息
		fmt.Println(err)
		return &pb.UserInfo{}, nil
	}
	// 返回结果
	return &pb.UserInfo{
		Name: u.Name,
		Age:  int32(u.Age),
		Sex:  pb.EnumSex(u.Sex),
	}, nil
}
