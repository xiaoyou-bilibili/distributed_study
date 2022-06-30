package test

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"server/ent"
	"server/ent/user"
	"testing"
)

func TestEntClient(t *testing.T) {
	// 建立数据库连接
	client, err := ent.Open("mysql", "root:xiaoyou@tcp(192.168.1.40:30000)/micro?parseTime=True")
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	// 自动创建表
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	//// 插入数据
	//save, err := client.User.Create().SetName("xiaoyou2").SetAge(10).SetSex(1).Save(context.Background())
	//if err != nil {
	//	fmt.Println("保存失败", err)
	//}
	//fmt.Println(save)

	// 测试查询
	u, err := client.User.
		Query().Where(user.ID(1)).
		Only(context.Background())

	fmt.Println(u.Name)
}
