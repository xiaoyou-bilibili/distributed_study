# Milvus 
## 搭建
```yaml
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: milvus
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: milvus
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: web
      k8s.kuboard.cn/name: milvus
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: web
        k8s.kuboard.cn/name: milvus
    spec:
      volumes:
        - name: volume-62miz
          persistentVolumeClaim:
            claimName: database
      containers:
        - name: etcd
          image: 'registry.xiaoyou66.com/quay.io/coreos/etcd:v3.5.0'
          command:
            - etcd
            - '-advertise-client-urls=http://127.0.0.1:2379'
            - '-listen-client-urls'
            - 'http://0.0.0.0:2379'
            - '--data-dir'
            - /etcd
          env:
            - name: ETCD_AUTO_COMPACTION_MODE
              value: revision
            - name: ETCD_AUTO_COMPACTION_RETENTION
              value: '1000'
            - name: ETCD_QUOTA_BACKEND_BYTES
              value: '4294967296'
          resources: {}
          volumeMounts:
            - name: volume-62miz
              mountPath: /etcd
              subPath: milvus/etcd
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
        - name: minio
          image: 'registry.xiaoyou66.com/minio/minio:RELEASE.2020-12-03T00-03-10Z'
          command:
            - minio
            - server
            - /minio_data
          env:
            - name: MINIO_ACCESS_KEY
              value: minioadmin
            - name: MINIO_SECRET_KEY
              value: minioadmin
          resources: {}
          volumeMounts:
            - name: volume-62miz
              mountPath: /minio_data
              subPath: milvus/minio
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
        - name: milvus
          image: 'registry.xiaoyou66.com/milvusdb/milvus:v2.0.2'
          command:
            - milvus
            - run
            - standalone
          env:
            - name: ETCD_ENDPOINTS
              value: '127.0.0.1:2379'
            - name: MINIO_ADDRESS
              value: '127.0.0.1:9000'
          resources: {}
          volumeMounts:
            - name: volume-62miz
              mountPath: /var/lib/milvus
              subPath: milvus/milvus
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: milvus
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
  name: milvus
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: milvus
spec:
  ports:
    - name: d4r2xm
      protocol: TCP
      port: 19530
      targetPort: 19530
      nodePort: 30510
  selector:
    k8s.kuboard.cn/layer: web
    k8s.kuboard.cn/name: milvus
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```
## 简单使用
```bash
# 安装python SDK
python3 -m pip install pymilvus==2.0.2
```

下面只简单介绍一下python SDK的使用
```python
from pymilvus import connections
# 创建一个链接
connections.connect(
	alias="default", 
	host='192.168.1.40', 
	port='30510'
)

# 创建一个schema
from pymilvus import CollectionSchema, FieldSchema, DataType
book_id = FieldSchema(
    name="book_id", 
    dtype=DataType.INT64, 
    is_primary=True, 
    )
word_count = FieldSchema(
    name="word_count", 
    dtype=DataType.INT64,  
    )
book_intro = FieldSchema(
    name="book_intro", 
    dtype=DataType.FLOAT_VECTOR, 
    dim=2
    )
schema = CollectionSchema(
    fields=[book_id, word_count, book_intro], 
    description="Test book search"
    )
collection_name = "book"
# 根据前面的schema来创建collection
from pymilvus import Collection
collection = Collection(
    name=collection_name, 
    schema=schema, 
    using='default', 
    shards_num=2,
    consistency_level="Strong"
    )


from pymilvus import utility
# 判断collection是否存在
a = utility.has_collection("book")
print(a)
# 下面我们查看一下详细信息
detail = collection = Collection("book")
print(detail)

# 列出所有的collection
utility.list_collections()

# 加载collection
collection = Collection("book")      # Get an existing collection.
collection.load()

# 下面准备一条数据
import random
data = [
    		[i for i in range(2000)],
		[i for i in range(10000, 12000)],
    		[[random.random() for _ in range(2)] for _ in range(2000)],
		]
# 插入这条数据
mr = collection.insert(data)


# 下面演示搜索数据
# 准备搜索的参数
search_params = {"metric_type": "L2", "params": {"nprobe": 10}}
# 进行向量搜索
results = collection.search(
	data=[[0.1, 0.2]], 
	anns_field="book_intro", 
	param=search_params, 
	limit=10, 
	expr=None,
	consistency_level="Strong"
)
print(results)
```

## 参考文档
- [官方文档](https://milvus.io/cn/docs/v2.0.x/overview.md)