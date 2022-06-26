# kratos

项目地址：https://github.com/go-kratos/kratos

## 准备工作

首先需要安装一些必要的依赖

```bash
# 需要安装一些protoc的编译器
sudo apt update && sudo apt install -y protobuf-compiler
```

## 使用指南

```bash
# 安装kratos
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
kratos upgrade
# 创建一个服务
kratos new demo1
# 拉取项目的依赖
go mod download
# 启动服务
kratos run
```

详细使用参考官方文档：https://go-kratos.dev/docs/intro/design

## 简单使用指南

kratos采用的是谷歌的proto3协议来构建的.

> 详细说明可以参考官方文档: https://developers.google.com/protocol-buffers/docs/proto3
> 中文的可以看一下这个：https://colobu.com/2019/10/03/protobuf-ultimate-tutorial-in-go/#%E5%8E%86%E5%8F%B2

我们可以使用官方提供的模板来构建服务，这里我们从零开始构建一个项目（项目名称叫demo1）。我们自己创建一个服务，然后自己新建一个`api`文件夹并新建 `demo.proto` 文件，文件内容如下

```proto
// 定义我们接口的版本
syntax = "proto3";
// 定义包名称
package api;

// 定义go安装包名称
option go_package = "demo1/api;demo";

// 定义我们的服务
service Demo {
    // 这里定义一个rpc服务，参数为HelloRequest，返回的结果是HelloReply
    rpc SayHello (HelloRequest) returns (HelloReply);
}

// 定义结构体信息
message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

然后我们根据这个接口定义来生成我们的代码

```bash
# 生成客户端代码
kratos proto client api/demo.proto
# 生成服务端代码-t表示生成的路径
kratos proto server api/demo.proto -t service
```

然后我们可以打开 `internal/service/demo.go` 里面就是服务端的处理代码了，我们直接返回请求的参数

> 注意生成的代码里面的方法大小写可能有有问题，需要自己改一下

```go
func (s *DemoService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{
        Message: "你好呀" + req.GetName(),
    }, nil
}
```

然后新建一个`main.go`来启动我们的服务

```go
package main

import (
    demo "demo1/api"
    "demo1/service"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/middleware/recovery"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "log"
)

func main() {
    s := service.NewDemoService()
    grpcSrv := grpc.NewServer(
        grpc.Address(":9000"),
        grpc.Middleware(recovery.Recovery()),
    )
    demo.RegisterDemoServer(grpcSrv, s)
    app := kratos.New(
        kratos.Name("demo"),
        kratos.Server(grpcSrv),
    )
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

然后我们写一个单元测试去连接一下这个服务

```go
import (
    "context"
    pb "demo1/api/demo"
    "fmt"
    "testing"

    "github.com/go-kratos/kratos/v2/transport/grpc"
)

func TestClient(t *testing.T) {
    conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:9000"))
    fmt.Println("连接", err)
    a := pb.NewDemoClient(conn)
    res, err := a.SayHello(context.Background(), &pb.HelloRequest{Name: "小游"})
    fmt.Println(err)
    fmt.Println("结果", res)
}
```

## 服务发现

这里使用`consul`来进行服务发现,这里只贴一下服务端和客户端的代码

```go
import (
    demo "demo1/api"
    "demo1/service"
    "github.com/go-kratos/kratos/contrib/registry/consul/v2"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/middleware/recovery"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/hashicorp/consul/api"
    "log"
)

func main() {
    s := service.NewDemoService()
    grpcSrv := grpc.NewServer(
        grpc.Address(":9000"),
        grpc.Middleware(recovery.Recovery()),
    )
    demo.RegisterDemoServer(grpcSrv, s)
    // 这里配置服务发现的地址
    config := api.DefaultConfig()
    // 这里配置成我们的consul的agent地址
    config.Address = "192.168.1.50:30571"
    consulClient, err := api.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    // 注册consul服务
    r := consul.New(consulClient)

    app := kratos.New(
        kratos.Name("kratos-demo"),
        kratos.Server(grpcSrv),
        kratos.Registrar(r),
    )
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

客户端代码

```go
import (
    "context"
    pb "demo1/api"
    "fmt"
    "github.com/go-kratos/kratos/contrib/registry/consul/v2"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/hashicorp/consul/api"
    "log"
    "testing"
)

func TestClient(t *testing.T) {
    config := api.DefaultConfig()
    config.Address = "192.168.1.50:30571"
    consulClient, err := api.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    // 注册consul服务
    r := consul.New(consulClient)
    conn, err := grpc.DialInsecure(
        context.Background(),
        grpc.WithEndpoint("discovery:///kratos-demo"),
        grpc.WithDiscovery(r))
    // 建立连接
    a := pb.NewDemoClient(conn)
    res, err := a.SayHello(context.Background(), &pb.HelloRequest{Name: "小游"})
    fmt.Println(err)
    fmt.Println("结果", res)
}
```