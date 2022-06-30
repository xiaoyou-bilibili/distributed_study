package main

import (
	pb "client/api"
	"client/loki"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
	"strconv"
)

var (
	// 使用 prometheus
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

func customMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		if tr, ok := transport.FromServerContext(ctx); ok {
			fmt.Println("operation:", tr.Operation())
		}
		reply, err = handler(ctx, req)
		return
	}
}

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
			semconv.ServiceNameKey.String("kratos-client"),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	// 自定义的日志初始化
	loki.Init()
	// 使用swagger
	openAPIhandler := openapiv2.NewHandler()
	// prometheus设置
	prometheus.MustRegister(_metricSeconds, _metricRequests)
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

	router := gin.Default()
	// 使用kratos中间件
	router.Use(kgin.Middlewares(recovery.Recovery(), customMiddleware))
	router.GET("/user/:id", func(ctx *gin.Context) {
		name := ctx.Param("id")
		loki.AppLog().Info("user is ", name)

		config := api.DefaultConfig()
		config.Address = "192.168.1.40:30571"
		consulClient, err := api.NewClient(config)
		if err != nil {
			log.Fatal(err)
		}

		// 获取配置
		conf, _, err := consulClient.KV().Get("hello", nil)

		fmt.Println("配置文件：", string(conf.Value))

		// 注册consul服务
		r := consul.New(consulClient)
		conn, err := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint("discovery:///kratos-server"),
			grpc.WithDiscovery(r),
			grpc.WithMiddleware(
				recovery.Recovery(),
				tracing.Client(),
				logging.Client(logger),
			))

		// 建立连接
		a := pb.NewUserClient(conn)
		id, err := strconv.Atoi(name)
		if err != nil {
			kgin.Error(ctx, errors.BadRequest("auth_error", "no authentication"))
			return
		}
		res, err := a.GetInfo(context.Background(), &pb.UserRequest{Id: int32(id)})

		if name == "error" {
			// 返回kratos error
			kgin.Error(ctx, errors.Unauthorized("auth_error", "no authentication"))
		} else {
			ctx.JSON(200, res)
		}
	})
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				tracing.Server(),
				logging.Server(logger),
				metrics.Server(
					metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
					metrics.WithRequests(prom.NewCounter(_metricRequests)),
				),
			),
		))
	httpSrv.HandlePrefix("/q/", openAPIhandler)
	httpSrv.Handle("/metrics", promhttp.Handler())
	// 主要路由设置
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("gin"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
