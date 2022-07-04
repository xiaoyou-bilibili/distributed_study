# Kvrocks

## 部署
```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: kvrocks
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: kvrocks
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: kvrocks
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: kvrocks
    spec:
      volumes:
        - name: volume-7mfcd
          persistentVolumeClaim:
            claimName: database
      containers:
        - name: kvrocks
          image: 'registry.xiaoyou66.com/kvrocks/kvrocks:latest'
          resources: {}
          volumeMounts:
            - name: volume-7mfcd
              mountPath: /tmp/kvrocks
              subPath: kvrocks
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: kvrocks
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
  name: kvrocks
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: kvrocks
spec:
  ports:
    - name: ynzcpk
      protocol: TCP
      port: 6666
      targetPort: 6666
      nodePort: 31669
  selector:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: kvrocks
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```