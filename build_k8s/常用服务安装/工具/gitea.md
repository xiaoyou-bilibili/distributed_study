### yaml
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    k8s.kuboard.cn/displayName: 私有git仓库
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
  name: gitea
  namespace: xiaoyou-tool
  resourceVersion: '733789'
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: gitea
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: gitea
    spec:
      containers:
        - image: registry.xiaoyou.com/gitea/gitea
          imagePullPolicy: Always
          name: gitea
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /data
              name: volume-fre7x
              subPath: gitea
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
        - name: volume-fre7x
          persistentVolumeClaim:
            claimName: tool-data
status:
  availableReplicas: 1
  conditions:
    - lastTransitionTime: '2022-05-22T12:01:54Z'
      lastUpdateTime: '2022-05-22T12:01:57Z'
      message: ReplicaSet "gitea-b87f47d" has successfully progressed.
      reason: NewReplicaSetAvailable
      status: 'True'
      type: Progressing
    - lastTransitionTime: '2022-05-22T12:21:52Z'
      lastUpdateTime: '2022-05-22T12:21:52Z'
      message: Deployment has minimum availability.
      reason: MinimumReplicasAvailable
      status: 'True'
      type: Available
  observedGeneration: 3
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1

---
apiVersion: v1
kind: Service
metadata:
  annotations: {}
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
  name: gitea
  namespace: xiaoyou-tool
  resourceVersion: '729295'
spec:
  clusterIP: 10.96.91.88
  clusterIPs:
    - 10.96.91.88
  internalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  ports:
    - name: ephnkk
      port: 3000
      protocol: TCP
      targetPort: 3000
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations: {}
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: gitea
  name: gitea
  namespace: xiaoyou-tool
  resourceVersion: '729452'
spec:
  ingressClassName: app
  rules:
    - host: gitea.xiaoyou.com
      http:
        paths:
          - backend:
              service:
                name: gitea
                port:
                  number: 3000
            path: /
            pathType: Prefix
status:
  loadBalancer:
    ingress:
      - ip: 192.168.1.52

```

## SSH推送
```bash
# 生成SSH私钥
ssh-keygen -t ed25519 -C "xiaoyou2333@foxmail.com"
# 显示pub后缀的内容
cat /home/xiaoyou/.ssh/id_ed25519.pub 
```