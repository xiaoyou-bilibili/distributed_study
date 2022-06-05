
```yaml
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: nextcloud
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: nextcloud
  annotations:
    k8s.kuboard.cn/displayName: 个人云盘
spec:
  replicas: 0
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: nextcloud
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: nextcloud
    spec:
      volumes:
        - name: volume-kramm
          persistentVolumeClaim:
            claimName: tool
        - name: data1
          nfs:
            server: 192.168.1.60
            path: /data/SD1
      containers:
        - name: nextcloud
          image: registry.xiaoyou.com/nextcloud
          ports:
            - containerPort: 80
              protocol: TCP
          resources: {}
          volumeMounts:
            - name: volume-kramm
              mountPath: /var/www/html
              subPath: nextcloud
            - name: data1
              mountPath: /data/SD1
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
  revisionHistoryLimit: 10
  progressDeadlineSeconds: 600

---
kind: Service
apiVersion: v1
metadata:
  name: nextcloud
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: nextcloud
spec:
  ports:
    - name: wybatw
      protocol: TCP
      port: 80
      targetPort: 80
      nodePort: 30016
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: nextcloud
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster

---
kind: Ingress
apiVersion: networking.k8s.io/v1
metadata:
  name: nextcloud
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: nextcloud
spec:
  ingressClassName: app
  rules:
    - host: cloud.xiaoyou.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: nextcloud
                port:
                  number: 80


```