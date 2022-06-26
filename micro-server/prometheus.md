## 部署

```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: prometheus
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: prometheus
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: monitor
      k8s.kuboard.cn/name: prometheus
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: monitor
        k8s.kuboard.cn/name: prometheus
    spec:
      volumes:
        - name: volume-ttjt4
          configMap:
            name: dev-conf
            items:
              - key: prometheus
                path: prometheus.yml
            defaultMode: 420
      containers:
        - name: prometheus
          image: 'registry.xiaoyou66.com/prom/prometheus:latest'
          resources: {}
          volumeMounts:
            - name: volume-ttjt4
              mountPath: /etc/prometheus/prometheus.yml
              subPath: prometheus.yml
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: prometheus
  podManagementPolicy: OrderedReady
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      partition: 0
  revisionHistoryLimit: 10

---
kind: Service
apiVersion: v1
metadata:
  name: prometheus
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: prometheus
spec:
  ports:
    - name: ajtcba
      protocol: TCP
      port: 9090
      targetPort: 9090
      nodePort: 30407
  selector:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: prometheus
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: dev-conf
  namespace: xiaoyou-dev
data:
  prometheus: >-
    # my global config

    global:
      scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
      evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
      # scrape_timeout is set to the global default (10s).

    # Alertmanager configuration

    alerting:
      alertmanagers:
        - static_configs:
            - targets:
              # - alertmanager:9093

    # Load rules once and periodically evaluate them according to the global
    'evaluation_interval'.

    rule_files:
      # - "first_rules.yml"
      # - "second_rules.yml"

    # A scrape configuration containing exactly one endpoint to scrape:

    # Here it's Prometheus itself.

    scrape_configs:
      # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
      - job_name: "prometheus"

        # metrics_path defaults to '/metrics'
        # scheme defaults to 'http'.

        static_configs:
          - targets: ["localhost:9090"]
      - job_name: 'kratos'
        static_configs:
          - targets: ['192.168.1.40:31715']
```

## kratos 接入

代码如下
```go
import (
    prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
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


// prometheus设置
prometheus.MustRegister(_metricSeconds, _metricRequests)
httpSrv := http.NewServer(
    http.Address(":8000"),
    http.Middleware(
        middleware.Chain(
            metrics.Server(
                metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
                metrics.WithRequests(prom.NewCounter(_metricRequests)),
            ),
        ),
    ))
httpSrv.Handle("/metrics", promhttp.Handler())
```
 
访问这个服务的 `/metrics` 就可以看到对应的服务了，然后自己到 `prometheus.yml` 把服务域名加上就行


