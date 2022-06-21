配置文件
```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: goland
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: goland
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: goland
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: goland
    spec:
      volumes:
        - name: volume-63jt2
          persistentVolumeClaim:
            claimName: tool
      containers:
        - name: goland
          image: registry.xiaoyou66.com/jetbrains/projector-goland
          resources: {}
          volumeMounts:
            - name: volume-63jt2
              mountPath: /home/projector-user
              subPath: idea/goland
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: goland
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
  name: goland
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: goland
spec:
  ports:
    - name: g7azgk
      protocol: TCP
      port: 8887
      targetPort: 8887
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: goland
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
  name: goland
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: goland
spec:
  ingressClassName: app
  rules:
    - host: goland.xiaoyou.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: goland
                port:
                  number: 8887
```