## 多集群部署
新手建议使用这个项目 https://github.com/labring/sealos

我们有5台机器，2个maser和3个node(都是ubuntu20.04)
- `192.168.1.30` master
- `192.168.1.31` master
- `192.168.1.32` node
- `192.168.1.33` node
- `192.168.1.34` node

> 这5台机器账号密码都一样，并且是允许root

```bash
# 使用root身份运行下面的内容
sealos init --user root --passwd xiaoyou --master 192.168.1.41 --master 192.168.1.42 --node 192.168.1.43 --node 192.168.1.44 --node 192.168.1.45 --pkg-url /home/xiaoyou/kube1.22.0.tar.gz --version v1.22.0
```

sealos init --user root --passwd xiaoyou --master 192.168.1.41 --pkg-url /home/xiaoyou/kube1.22.0.tar.gz --version v1.22.0

其他常用命令

```bash
# 添加master
sealos join --master 192.168.0.6 --master 192.168.0.7
# 添加节点
sealos join --node 192.168.0.6 --node 192.168.0.7
# 删除某个节点
sealos clean --master 192.168.0.2
sealos clean --node 192.168.0.5
# 删除所有集群
sealos clean --all
# 查看所有节点
kubectl get nodes
# 开启端口服务
kubectl proxy --address='0.0.0.0'  --accept-hosts='^*$' --port=8081
```

token访问
https://blog.csdn.net/weixin_38320674/article/details/107328982

https://kuboard.cn/install/v3/install-in-k8s.html#%E5%AE%89%E8%A3%85