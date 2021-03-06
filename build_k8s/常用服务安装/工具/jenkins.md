参考：https://github.com/jenkinsci/docker/blob/master/README.md

```yaml

---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: jenkins
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: jenkins
  annotations:
    k8s.kuboard.cn/displayName: 自动化部署
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: jenkins
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: jenkins
    spec:
      volumes:
        - name: volume-eyf72
          persistentVolumeClaim:
            claimName: tool
      containers:
        - name: jenkins
          image: 'registry.xiaoyou.com/jenkins/jenkins:lts-jdk11'
          resources: {}
          volumeMounts:
            - name: volume-eyf72
              mountPath: /var/jenkins_home
              subPath: jenkins
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
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
  name: jellyfin
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: jellyfin
spec:
  ports:
    - name: ed5k6b
      protocol: TCP
      port: 8096
      targetPort: 8096
      nodePort: 30020
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: jellyfin
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
  name: jenkins
  namespace: xiaoyou-tool
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: jenkins
spec:
  ingressClassName: app
  rules:
    - host: jenkins.xiaoyou.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: jenkins
                port:
                  number: 8080


```