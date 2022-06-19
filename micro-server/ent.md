这里我们使用facebook开源的数据库操作框架

文档地址：https://entgo.io/zh/docs/getting-started

## 简单使用
```bash
# 安装工具
go get entgo.io/ent/cmd/ent
# 初始化我们的实体
ent init User
# 也可以用下面这个命令
go run entgo.io/ent/cmd/ent init User

# 初始化好的实体会在 project/ent/schema/ 目录下
# 我们可以自己去配置字段

# 配置好之后可以使用下面这个命令自动生成相关的增删改查文件
go generate ./ent
```

下面演示一个最简单的插入和查询，其他的自行参考官方文档

```go
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
	client, err := ent.Open("mysql", "root:xiaoyou@tcp(192.168.1.50:30000)/micro?parseTime=True")
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	// 自动创建表
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	
	// 插入数据
	save, err := client.User.Create().SetName("xiaoyou2").SetAge(10).SetSex(1).Save(context.Background())
	if err != nil {
		fmt.Println("保存失败", err)
	}
	fmt.Println(save)

	// 测试查询
	u, err := client.User.
		Query().Where(user.ID(1)).
		Only(context.Background())

	fmt.Println(u.Name)
}
```