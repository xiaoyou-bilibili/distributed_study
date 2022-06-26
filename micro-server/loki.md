主要要搭建服务
```bash
sudo docker pull grafana/loki
sudo docker pull grafana/promtail
```

## 部署

```yaml
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: loki
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: loki
  annotations: {}
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: monitor
      k8s.kuboard.cn/name: loki
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: monitor
        k8s.kuboard.cn/name: loki
    spec:
      volumes:
        - name: volume-pwybn
          configMap:
            name: dev-conf
            items:
              - key: promtail
                path: config.yml
            defaultMode: 420
      containers:
        - name: loki
          image: 'registry.xiaoyou66.com/grafana/loki:latest'
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
        - name: promtail
          image: 'registry.xiaoyou66.com/grafana/promtail:latest'
          resources: {}
          volumeMounts:
            - name: volume-pwybn
              mountPath: /etc/promtail/config.yml
              subPath: config.yml
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 25%
  revisionHistoryLimit: 10
  progressDeadlineSeconds: 600

---
kind: Service
apiVersion: v1
metadata:
  name: loki
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: loki
spec:
  ports:
    - name: ixrkcd
      protocol: TCP
      port: 3100
      targetPort: 3100
      nodePort: 30913
  selector:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: loki
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
  promtail: |
    server:
      http_listen_port: 9080
      grpc_listen_port: 0

    positions:
      filename: /tmp/positions.yaml

    clients:
      - url: http://127.0.0.1:3100/loki/api/v1/push

    scrape_configs:
    - job_name: system
      static_configs:
      - targets:
          - localhost
        labels:
          job: varlogs
          __path__: /var/log/*log
```


## golang接入

实际上这个日志应该是通过promtail来收集的，这里主要是以演示功能为主
```go
package loki

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// UploadLog 上报日志的格式
type UploadLog struct {
	Streams []Stream `json:"streams"`
}

// 自己实现一个loki的钩子，自己实现一下go的接口就可以了
type lokiHook struct{}

// Levels 这里表示我们过滤所有日志
func (h *lokiHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire 这个是钩子处理
func (h *lokiHook) Fire(e *log.Entry) error {
	// 这里我们打两个标签
	data := UploadLog{Streams: []Stream{{
		Stream: map[string]string{
			"level": e.Level.String(),
			"app":   e.Data["app"].(string),
		},
		Values: [][]string{{
			strconv.FormatInt(time.Now().UnixNano(), 10),
			e.Message,
		}},
	}}}
	a, _ := json.Marshal(e)
	fmt.Print(string(a))
	// 直接发送，这里不管是否成功
	_, _ = HttpPostJson("http://192.168.1.40:30913/loki/api/v1/push", data)
	return nil
}

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	// 日志直接打印到标准输出，不保存到本地
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// 添加钩子，通过钩子来进行数据上报
	log.AddHook(&lokiHook{})
}

// AppLog 打印APP日志
func AppLog() *log.Entry {
	return log.WithField("app", "kratos")
}


// HttpPostJson 发送json格式的数据
func HttpPostJson(url string, data interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("解析JSON数据失败")
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, errors.New("发送请求失败")
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("解析请求体失败")
	}
	return s, nil
}
```

使用的时候可以先初始化，然后直接使用来打印日志了
```go
loki.AppLog().Info("user is ", name)
```