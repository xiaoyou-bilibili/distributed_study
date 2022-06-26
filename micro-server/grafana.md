## 搭建

```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: grafana
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: grafana
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: monitor
      k8s.kuboard.cn/name: grafana
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: monitor
        k8s.kuboard.cn/name: grafana
    spec:
      containers:
        - name: grafana
          image: 'registry.xiaoyou66.com/grafana/grafana-oss:9.0.1'
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: grafana
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
  name: grafana
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: grafana
spec:
  ports:
    - name: pwwkxp
      protocol: TCP
      port: 3000
      targetPort: 3000
  selector:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: grafana
  type: ClusterIP
  sessionAffinity: None
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster

---
kind: Ingress
apiVersion: networking.k8s.io/v1
metadata:
  name: grafana
  namespace: xiaoyou-dev
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: grafana
spec:
  ingressClassName: app
  rules:
    - host: grafana.xiaoyou.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: grafana
                port:
                  number: 3000


```


