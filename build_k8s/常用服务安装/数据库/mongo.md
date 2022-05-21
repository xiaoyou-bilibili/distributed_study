## mongoDB 
地址：https://hub.docker.com/_/mongo


## yaml文件
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    k8s.kuboard.cn/displayName: MongoDB
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mongo
  name: mongo
  namespace: xiaoyou-database
  resourceVersion: '416158'
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 5
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: mongo
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: mongo
    spec:
      containers:
        - args:
            - '-f'
            - /etc/mongod.conf
          command:
            - mongod
          image: registry.xiaoyou.com/mongo
          imagePullPolicy: Always
          name: mongo
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /data/lib/mongodb
              name: volume-i3mpf
              subPath: mongo
            - mountPath: /etc/mongod.conf
              name: volume-nzab6
              subPath: mongod.conf
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
        - name: volume-i3mpf
          persistentVolumeClaim:
            claimName: database-storage
        - configMap:
            defaultMode: 420
            items:
              - key: mongo
                path: mongod.conf
            name: database-conf
          name: volume-nzab6
status:
  availableReplicas: 1
  collisionCount: 1
  conditions:
    - lastTransitionTime: '2022-05-21T12:25:38Z'
      lastUpdateTime: '2022-05-21T12:25:38Z'
      message: Deployment has minimum availability.
      reason: MinimumReplicasAvailable
      status: 'True'
      type: Available
    - lastTransitionTime: '2022-05-20T11:11:40Z'
      lastUpdateTime: '2022-05-21T12:25:38Z'
      message: ReplicaSet "mongo-cc8d997bd" has successfully progressed.
      reason: NewReplicaSetAvailable
      status: 'True'
      type: Progressing
  observedGeneration: 59
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1


```

mongo的配置文件如下

```bash
# mongod.conf

# for documentation of all options, see:
#   http://docs.mongodb.org/manual/reference/configuration-options/

# Where and how to store data.
storage:
  dbPath: /data/lib/mongodb
  journal:
    enabled: true
#  engine:
#  wiredTiger:

# where to write logging data.
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongod.log

# network interfaces
net:
  port: 27017
  bindIp: 0.0.0.0

```