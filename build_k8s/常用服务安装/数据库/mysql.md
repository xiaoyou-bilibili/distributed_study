## mysql
官方镜像地址：https://hub.docker.com/_/mysql


- 密码通过`MYSQL_ROOT_PASSWORD`来配置
- 存储路径在 `/var/lib/mysql`

## yaml文件
```bash

---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: mysql-8-0
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mysql-8-0
  annotations:
    k8s.kuboard.cn/displayName: MySQL8.0
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: mysql-8-0
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/layer: db
        k8s.kuboard.cn/name: mysql-8-0
    spec:
      volumes:
        - name: volume-dez7x
          persistentVolumeClaim:
            claimName: database
      containers:
        - name: mysql
          image: 'registry.xiaoyou66.com/mysql:8.0'
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: xiaoyou
          resources: {}
          volumeMounts:
            - name: volume-dez7x
              mountPath: /var/lib/mysql
              subPath: mysql8-0
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      schedulerName: default-scheduler
  serviceName: mysql-8-0
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
  name: mysql-8-0
  namespace: xiaoyou-database
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mysql-8-0
spec:
  ports:
    - name: rp5aam
      protocol: TCP
      port: 3306
      targetPort: 3306
      nodePort: 30000
  selector:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mysql-8-0
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```

### mysql新建用户

docker项目最好单独新开一个用户，避免不相关数据表被修改

参考：https://www.cnblogs.com/chanshuyi/p/mysql_user_mng.html

```bash
# 新建用户名和密码都为xiaoyou
create user xiaoyou identified by 'xiaoyou';
# 给xiaoyou授予demo数据库的权限
grant all privileges on deno.* to xiaoyou@'%';
# 刷新权限
flush privileges;
# 查看用户的权限 
show grants for 'xiaoyou';
```