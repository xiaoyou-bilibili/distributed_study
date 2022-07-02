# MySQL数据库

## 搭建

```yaml
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

## 常用命令

```bash
# 安装mysql客户端
sudo apt install mysql-client-core-8.0
# 连接数据库 -h 表示主机名 -P表示端口 -u表示用户 -p可以指定密码
mysql -h 192.168.1.40 -P 30001 -u root -p
# 修改用户密码
mysqladmin -u 用户名 -p 旧密码 password 新密码
# 显示数据库
show databases;
# 创建一个demo数据库
create database demo;
# 删除demo数据库
drop database demo;
# 使用数据库
use demo
# 显示所有的表
show tables;

#创建一个表
create table test(
    id int(4) not null primary key auto_increment,
    name char(20) not null,
    age int not null default 1);

# 我们查看表的结构
desc test;
# 删除表
drop table test
# 修改表名
rename table test to demo;
# 修改表结构
alter table demo add address char(255) default '';
# 添加数据
insert into demo (name,age,address) values ('小游',10,'上海');
# 删除数据
delete from demo where id = 1;
# 修改数据
update demo set name='小白',age=18 where id =2;
# 查询数据
select * from demo;
# 数据库备份 -p 表示数据库名，后面跟着表名字，然后a.sql 就是导出结果
mysqldump -h 192.168.1.40 -P 30001 -u root -p demo demo > a.sql
```

## 可视化管理工具

- navicate

- omnidb

```bash
sudo docker run -itd --name db -v /xiaoyou/db:/home/omnidb/.omnidb/omnidb-server -p 8027:8000 omnidbteam/omnidb:latest
```

## 参考

- [MySQL命令大全](http://c.biancheng.net/cpp/html/1446.html)
