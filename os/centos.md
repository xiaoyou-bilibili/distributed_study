## CentOS

官网：https://www.centos.org/


## 常见问题

### 网络设置

```bash
# 使用ip addr查看系统网络名称
# 修改网络配置
vi /etc/sysconfig/network-scripts/ifcfg-<系统网络名称>
#将ONBOOT的值改为yes
```

### 给普通用户添加sudo权限
```bash
# 切换到root用户
visudo
# 找到下面这一行
root  ALL=(ALL)    ALL
# 自己按照这个格式在后面添加(xiaoyou就是自己的用户名)
xiaoyou  ALL=(ALL)    ALL
```