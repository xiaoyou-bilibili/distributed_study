## 搜索引擎
参考： https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html

```bash
sudo docker pull docker.elastic.co/elasticsearch/elasticsearch:8.2.0
sudo docker pull docker.elastic.co/kibana/kibana:8.2.0
```

### es

```bash
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  annotations:
    k8s.kuboard.cn/displayName: 搜索引擎
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: elasticsearch
  name: elasticsearch
  namespace: xiaoyou-database
  resourceVersion: '445313'
spec:
  podManagementPolicy: OrderedReady
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: elasticsearch
  serviceName: elasticsearch
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: elasticsearch
    spec:
      containers:
        - env:
            - name: discovery.type
              value: single-node
            - name: ES_JAVA_OPTS
              value: '-Xms200m -Xmx200m'
          image: >-
            registry.xiaoyou.com/docker.elastic.co/elasticsearch/elasticsearch:8.2.0
          imagePullPolicy: Always
          name: elasticsearch
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /usr/share/elasticsearch/data
              name: volume-3r8ch
              subPath: elasticsearch/data
            - mountPath: /usr/share/elasticsearch/plugins
              name: volume-3r8ch
              subPath: elasticsearch/plugins
            - mountPath: /usr/share/elasticsearch/config/elasticsearch.yml
              name: volume-eh266
              subPath: elasticsearch.yml
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
        - name: volume-3r8ch
          persistentVolumeClaim:
            claimName: database-storage
        - configMap:
            defaultMode: 420
            items:
              - key: elasticsearch
                path: elasticsearch.yml
            name: database-conf
          name: volume-eh266
  updateStrategy:
    rollingUpdate:
      partition: 0
    type: RollingUpdate
status:
  availableReplicas: 1
  collisionCount: 0
  currentReplicas: 1
  currentRevision: elasticsearch-54c56887fd
  observedGeneration: 17
  readyReplicas: 1
  replicas: 1
  updateRevision: elasticsearch-54c56887fd
  updatedReplicas: 1
```

### kibana

```bash
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    k8s.kuboard.cn/displayName: ES管理
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: kibana
  name: kibana
  namespace: xiaoyou-database
  resourceVersion: '444773'
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: kibana
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
        k8s.kuboard.cn/name: kibana
    spec:
      containers:
        - env:
            - name: ELASTICSEARCH_HOSTS
              value: 'http://192.168.1.50:30003'
          image: 'registry.xiaoyou.com/docker.elastic.co/kibana/kibana:8.2.0'
          imagePullPolicy: IfNotPresent
          name: kibana
          ports:
            - containerPort: 5601
              protocol: TCP
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 1
  conditions:
    - lastTransitionTime: '2022-05-21T00:47:18Z'
      lastUpdateTime: '2022-05-21T14:30:47Z'
      message: ReplicaSet "kibana-7dc9455c88" has successfully progressed.
      reason: NewReplicaSetAvailable
      status: 'True'
      type: Progressing
    - lastTransitionTime: '2022-05-21T14:31:55Z'
      lastUpdateTime: '2022-05-21T14:31:55Z'
      message: Deployment has minimum availability.
      reason: MinimumReplicasAvailable
      status: 'True'
      type: Available
  observedGeneration: 7
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1


```


配置文件如下，配置文件在 `/usr/share/elasticsearch/config/elasticsearch.yml` 中
```
cluster.name: "docker-cluster"
network.host: 0.0.0.0
# 关闭密码
xpack.security.enabled: false
```
