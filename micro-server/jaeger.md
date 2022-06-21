Jaeger 是由Uber 开源的分布式追踪系统，它采用Go语言编写，主要借鉴了 Google Dapper 论文和 Zipkin 的设计，兼容 OpenTracing 以及 Zipkin 追踪格式，目前已成为CNCF基金会的开源项目。


部署可以参考：https://www.jaegertracing.io/docs/1.35/

## 单机部署

```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: jaeger
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: jaeger
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: monitor
      k8s.kuboard.cn/name: jaeger
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: monitor
        k8s.kuboard.cn/name: jaeger
    spec:
      containers:
        - name: jaeger
          image: registry.xiaoyou66.com/jaegertracing/all-in-one
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: jaeger
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
  name: jaeger
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: jaeger
spec:
  ports:
    - name: bc7epy
      protocol: TCP
      port: 16686
      targetPort: 16686
      nodePort: 31708
    - name: cz7brb
      protocol: TCP
      port: 14268
      targetPort: 14268
      nodePort: 30878
  selector:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: jaeger
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```