# InfluxDB 

## 部署
```bash
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: influx
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: influx
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: influx
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: influx
    spec:
      volumes:
        - name: volume-4b4j5
          persistentVolumeClaim:
            claimName: database
      containers:
        - name: influx
          image: 'registry.xiaoyou66.com/library/influxdb:latest'
          resources: {}
          volumeMounts:
            - name: volume-4b4j5
              mountPath: /var/lib/influxdb
              subPath: influxdb
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: influx
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
  name: influx
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: influx
spec:
  ports:
    - name: jwengc
      protocol: TCP
      port: 8086
      targetPort: 8086
      nodePort: 32206
  selector:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: influx
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```

> 搭建好后直接访问8086端口就可以了

## 常用命令

这里1.0和2.0两个版本相差巨大，所以这里后续有需要再进行补充

```bash

```

## 参考文档
- https://jasper-zhang1.gitbooks.io/influxdb/content/Guide/writing_data.html
- https://www.jianshu.com/p/d2935e99006e