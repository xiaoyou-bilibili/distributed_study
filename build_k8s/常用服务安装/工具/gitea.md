### yaml
```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: gitea
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
  annotations:
    k8s.kuboard.cn/displayName: 私有git服务
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: gitea
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: gitea
    spec:
      volumes:
        - name: volume-3kkce
          persistentVolumeClaim:
            claimName: tool
      containers:
        - name: gitea
          image: registry.xiaoyou66.com/gitea/gitea
          resources: {}
          volumeMounts:
            - name: volume-3kkce
              mountPath: /data
              subPath: gitea
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: gitea
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
  name: gitea
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
spec:
  ports:
    - name: ebep8f
      protocol: TCP
      port: 3000
      targetPort: 3000
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
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
  name: gitea
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
spec:
  ingressClassName: app
  rules:
    - host: gitea.xiaoyou.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: gitea
                port:
                  number: 3000
```

## SSH推送
```bash
# 生成SSH私钥
ssh-keygen -t ed25519 -C "xiaoyou2333@foxmail.com"
# 显示pub后缀的内容
cat /home/xiaoyou/.ssh/id_ed25519.pub 
```