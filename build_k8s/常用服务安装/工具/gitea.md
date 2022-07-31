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
ssh-keygen -t rsa -C "xiaoyou2333@foxmail.com"
# 显示pub后缀的内容,然后到仓库里面设置一下
cat /home/xiaoyou/.ssh/id_rsa.pub
git remote add origin ssh://git@192.168.1.40:32045/index/music-player.git
git push -u origin master
```

自己去修改一下`app.ini`文件
```conf
[server]
APP_DATA_PATH    = /data/gitea
DOMAIN           = 192.168.1.40
SSH_DOMAIN       = 192.168.1.40
HTTP_PORT        = 3000
ROOT_URL         = https://git.xiaoyou.host/
DISABLE_SSH      = false
SSH_PORT         = 32045
SSH_LISTEN_PORT  = 22
LFS_START_SERVER = true
LFS_CONTENT_PATH = /data/git/lfs
LFS_JWT_SECRET   = ajP-ErA4aHjM_cEp9ri98DO7iX_g2aSKEy6aeqiwx4I
OFFLINE_MODE     = false
```