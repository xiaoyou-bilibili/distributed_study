package main

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
	user "server/api"
	"server/service"
)

var serverName = "kratos-server"

func setTracerProvider(url string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serverName),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	// 配置链路追踪
	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "trace_id", tracing.TraceID())
	logger = log.With(logger, "span_id", tracing.SpanID())
	log := log.NewHelper(logger)

	url := "http://192.168.1.40:30878/api/traces"
	if os.Getenv("jaeger_url") != "" {
		url = os.Getenv("jaeger_url")
	}
	err := setTracerProvider(url)
	if err != nil {
		log.Error(err)
	}

	s := service.NewUserService()
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				tracing.Server(),
				logging.Server(logger),
			),
		),
	)
	user.RegisterUserServer(grpcSrv, s)
	// 这里配置服务发现的地址
	config := api.DefaultConfig()
	// 这里配置成我们的consul的agent地址
	config.Address = "192.168.1.40:30571"
	consulClient, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// 注册consul服务
	r := consul.New(consulClient)

	app := kratos.New(
		kratos.Name(serverName),
		kratos.Server(grpcSrv),
		kratos.Registrar(r),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
