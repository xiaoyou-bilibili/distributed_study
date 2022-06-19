```yaml

---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: sftp
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: sftp
  annotations:
    k8s.kuboard.cn/displayName: ftp服务
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: sftp
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: sftp
    spec:
      volumes:
        - name: volume-bwkxw
          persistentVolumeClaim:
            claimName: database
      containers:
        - name: sftp
          image: 'registry.xiaoyou66.com/atmoz/sftp:debian'
          args:
            - 'xiaoyou:xiaoyou:::upload'
          resources: {}
          volumeMounts:
            - name: volume-bwkxw
              mountPath: /home/xiaoyou/upload
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
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
  revisionHistoryLimit: 1
  progressDeadlineSeconds: 600

---
kind: Service
apiVersion: v1
metadata:
  name: sftp
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: sftp
spec:
  ports:
    - name: mbsknp
      protocol: TCP
      port: 22
      targetPort: 22
      nodePort: 30189
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: sftp
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```