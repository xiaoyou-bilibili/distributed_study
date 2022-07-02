# ClickHouse

> 这个clickhouse 给我的感觉很像mysql，这里就先不管了，后面我再详细研究

## 部署

```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: clickhouse
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: clickhouse
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: clickhouse
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: clickhouse
    spec:
      volumes:
        - name: volume-dr4jm
          persistentVolumeClaim:
            claimName: database
      containers:
        - name: clickhouse
          image: 'registry.xiaoyou66.com/clickhouse/clickhouse-server:latest'
          resources: {}
          volumeMounts:
            - name: volume-dr4jm
              mountPath: /var/lib/clickhouse
              subPath: clickhouse
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: clickhouse
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
  name: clickhouse
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: clickhouse
spec:
  ports:
    - name: wbrym7
      protocol: TCP
      port: 8123
      targetPort: 8123
      nodePort: 32380
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: clickhouse
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```

## 常用命令
```bash
# 安装连接工具
sudo apt install clickhouse-client
# 连接远程(注意，TCP的需要开放9000端口)
clickhouse-client -h 192.168.1.40 --port 32380
# 显示所有的数据库
show databases
# 创建数据库
CREATE DATABASE IF NOT EXISTS tutorial
# 我们创建一个表
create table if not exists order
(
  id Int64 COMMENT '订单id', 
  datetime DateTime COMMENT '订单日期',
  name String COMMENT '商品名称',
  price Decimal32(2) COMMENT '商品价格',
  user_id Int64 COMMENT '用户id'
) engine = MergeTree 
partition by toYYYYMM(datetime)
order by id 
# 插入一条测试数据
insert into order (id,datetime,name,price,user_id) values (1, '2021-03-09 21:42:00', '大力丸', 999.99, 202003090001)

# 查询数据
select * from order order by id desc;
```


## GUI连接工具
- https://github.com/tabixio/tabix

## 参考文档
- [官方文档](https://clickhouse.com/docs/zh/)