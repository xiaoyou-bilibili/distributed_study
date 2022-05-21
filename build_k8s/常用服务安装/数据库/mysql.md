## mysql
官方镜像地址：https://hub.docker.com/_/mysql


- 密码通过`MYSQL_ROOT_PASSWORD`来配置
- 存储路径在 `/var/lib/mysql`

## yaml文件
```bash
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    k8s.kuboard.cn/displayName: MySQL8.0
  labels:
    k8s.kuboard.cn/layer: db
    k8s.kuboard.cn/name: mysql
  name: mysql
  namespace: xiaoyou-database
  resourceVersion: '414519'
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 5
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: db
      k8s.kuboard.cn/name: mysql
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
        k8s.kuboard.cn/name: mysql
    spec:
      containers:
        - env:
            - name: MYSQL_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: passwd
                  name: database-conf
          image: 'registry.xiaoyou.com/mysql:8.0'
          imagePullPolicy: Always
          name: mysql
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /var/lib/mysql
              name: volume-t8jzk
              subPath: mysql8-0
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
        - name: volume-t8jzk
          persistentVolumeClaim:
            claimName: database-storage
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