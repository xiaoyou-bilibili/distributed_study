# MongoDB

## 部署
```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: mongo
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mongo
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: mongo
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: mongo
    spec:
      volumes:
        - name: volume-zx7bb
          persistentVolumeClaim:
            claimName: database
        - name: volume-hfwri
          configMap:
            name: database-conf
            items:
              - key: mongo
                path: mongod.conf
            defaultMode: 420
      containers:
        - name: mongo
          image: 'registry.xiaoyou66.com/library/mongo:latest'
          command:
            - mongod
          args:
            - '-f'
            - /etc/mongod.conf
          resources: {}
          volumeMounts:
            - name: volume-zx7bb
              mountPath: /data/lib/mongodb
              subPath: mongo
            - name: volume-hfwri
              mountPath: /etc/mongod.conf
              subPath: mongod.conf
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: mongo
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
  name: mongo
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mongo
spec:
  ports:
    - name: ckckwy
      protocol: TCP
      port: 27017
      targetPort: 27017
      nodePort: 30115
  selector:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mongo
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```

mongo的配置文件如下
```yaml
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

## 常用命令

```bash
# 安装连接工具
sudo apt install mongodb-clients
# 连接mongo
mongo mongodb://user:pass@localhost:port/
# mongodb://192.168.1.40:30115
# 查看所有的数据库
show dbs

# 使用某个数据库(如果不存在则会自己新建)
use demo
# 查看当前使用的数据库
db
# 删除数据库(需要先切换到某个数据库)
db.dropDatabase()
# 如果想删除集合可以这样
db.collection.drop()
# 查看所有的集合，下面两个命令是一样的
show tables 
show collections

# 创建集合
db.createCollection("test")


# 插入文档
db.test.insert({"name":"xiaoyou"})
# 查看文档
db.test.find()
# 更新文档
db.test.update({"name":"xiaoyou"},{$set:{"name":"xiaobai"}})
# 删除文档
db.test.remove({"name":"xiaobai"})
```

## 可视化工具

- navicate for mongo
- Studio 3T


## 参考文章

- https://www.runoob.com/mongodb/mongodb-tutorial.html