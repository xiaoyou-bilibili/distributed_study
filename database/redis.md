# Redis

## 搭建

```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: redis
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: redis
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: redis
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: redis
    spec:
      containers:
        - name: redis
          image: 'registry.xiaoyou66.com/library/redis:latest'
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: redis
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
  name: redis
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: redis
spec:
  ports:
    - name: mm77zz
      protocol: TCP
      port: 6379
      targetPort: 6379
      nodePort: 31686
  selector:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: redis
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```

## 常用命令

``` bash
# 安装redis客户端
sudo apt install redis-tools
# 连接远程redis
redis-cli -h host -p port -a password
# 测试连接
PING
# 设置一个键，值为xiaoyou
SET name xiaoyou
# 获取一个键
GET name
# 设置key的过期时间，按秒计算
expire name 10
# 删除一个键
DEL name

# 设置一个hash(其实就是一个map 这里设置了两个key 第一个是name第二个是description)
hmset maps name "xiaoyou" description "66"
# 获取所有的内容 
hgetall maps
# 获取map里面的某个指定的字段
hget maps name

# 给列表添加元素(这里给names 添加了两个元素)
lpush names xiaoyou xiaobai
# 遍历我们的列表,分别为名字，开始位置和结束位置
lrange names 0 10


# 我们给集合添加元素（集合不能出现重复的元素）
sadd database mongo redis mysql
# 获取集合的元素
smembers database

# 有序集合，有序集合有分数，我们可以给前面设置分数，后面就会按照从小到大的分数进行排序
zadd app 1 bilibili
# 获取有序集合的结果，这里会把分数给输出
zrange app 0 10 withscores
```



## 可视化管理工具

- [AnotherRedisDesktopManager](https://github.com/qishibo/AnotherRedisDesktopManager)

## 参考文章

- https://www.runoob.com/redis/redis-tutorial.html