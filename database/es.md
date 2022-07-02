# Elasticsearch
## 搭建
```yaml

---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: elasticsearch
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: elasticsearch
  annotations:
    k8s.kuboard.cn/displayName: 搜索引擎
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: elasticsearch
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: elasticsearch
    spec:
      volumes:
        - name: volume-3r8ch
          persistentVolumeClaim:
            claimName: database
        - name: volume-eh266
          configMap:
            name: database-conf
            items:
              - key: elasticsearch
                path: elasticsearch.yml
            defaultMode: 420
      containers:
        - name: elasticsearch
          image: >-
            registry.xiaoyou66.com/docker.elastic.co/elasticsearch/elasticsearch:8.2.0
          env:
            - name: discovery.type
              value: single-node
            - name: ES_JAVA_OPTS
              value: '-Xms200m -Xmx200m'
          resources: {}
          volumeMounts:
            - name: volume-3r8ch
              mountPath: /usr/share/elasticsearch/data
              subPath: elasticsearch/data
            - name: volume-3r8ch
              mountPath: /usr/share/elasticsearch/plugins
              subPath: elasticsearch/plugins
            - name: volume-eh266
              mountPath: /usr/share/elasticsearch/config/elasticsearch.yml
              subPath: elasticsearch.yml
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
        - name: kibana
          image: 'registry.xiaoyou66.com/docker.elastic.co/kibana/kibana:8.2.0'
          env:
            - name: ELASTICSEARCH_HOSTS
              value: 'http://127.0.0.1:9200'
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: elasticsearch
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
  name: elasticsearch
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: elasticsearch
spec:
  ports:
    - name: es-api
      protocol: TCP
      port: 9200
      targetPort: 9200
      nodePort: 31981
    - name: hkmkmm
      protocol: TCP
      port: 9300
      targetPort: 9300
      nodePort: 32261
    - name: kibana
      protocol: TCP
      port: 5601
      targetPort: 5601
      nodePort: 32621
  selector:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: elasticsearch
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
  name: elasticsearch
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: elasticsearch
spec:
  ingressClassName: app
  rules:
    - host: kibana.xiaoyou.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: elasticsearch
                port:
                  number: 5601
```

配置文件内容如下
```yaml
cluster.name: "docker-cluster"
network.host: 0.0.0.0
# 关闭密码
xpack.security.enabled: false
```