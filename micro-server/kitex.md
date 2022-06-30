# kitex

kitex是字节开源的微服务框架

代码地址：https://github.com/cloudwego/kitex

官方文档：https://www.cloudwego.io/zh/docs/kitex/getting-started/

## 环境搭建
```bash
# 安装kitex
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
# 安装thriftgo
go install github.com/cloudwego/thriftgo@latest
```

## 简单demo

我们先写一个 `hello.thrift` 文件，内容如下：
```thrift
namespace go api

struct Request {
    1: string message
}

struct Response {
    1: string message
}

service Hello {
    Response echo(1: Request req)
}
```

然后我们根据thrift来生成一下代码

```bash
kitex -module "kitex" -service kitex_demo api/hello.thrift
```

具体的服务代码会在 `handler.go` 里面，自己重写一下就可以了

> 注意，thrift的版本不能太高了，要不然无法运行，如果无法运行可以替换成下面的版本

```
github.com/apache/thrift v0.13.0
```


最后我们写一个测试代码去连接这个服务，代码如下
```go
import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"kitex/kitex_gen/api"
	"kitex/kitex_gen/api/hello"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	c, err := hello.NewClient("kitex_demo", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}

	res,err:= c.Echo(context.Background(),&api.Request{Message: "小游"})
	fmt.Println(res.Message)
}
```
