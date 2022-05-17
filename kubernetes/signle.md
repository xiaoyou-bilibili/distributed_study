## 单机搭建
参考：https://kubernetes.io/zh/docs/tasks/tools/

## 安装
安装 kubectl，参考：https://kubernetes.io/zh/docs/tasks/tools/install-kubectl-linux/
```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

安装minikube，安装参考：https://minikube.sigs.k8s.io/docs/start/

```bash
# 首先安装docker，可以参考 https://docs.docker.com/engine/install/debian/
# 守护进程启动
systemctl start docker
sudo systemctl daemon-reload
systemctl enable docker
systemctl restart docker
# docker添加用户权限
sudo usermod -aG docker $USER && newgrp docker
# 安装minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
# 查看自己的版本
minikube version
# 启动一个集群
minikube start
# 我们查看一下创建的集群信息
docker ps 
```

查看集群信息
```bash
# 查看集群信息
kubectl cluster-info
# 获取节点
kubectl get nodes
# 返回下面内容说明安装完毕
NAME       STATUS   ROLES                  AGE   VERSION
minikube   Ready    control-plane,master   93s   v1.23.3
```

### 部署应用

```bash
# 先确保自己安装了kubectl
kubectl version
# 并且有节点
kubectl get nodes
# 下面安装一个最简单的应用
kubectl create deployment kubernetes-bootcamp --image=gcr.io/google-samples/kubernetes-bootcamp:v1
# 查看我们部署的应用
kubectl get deployments
# 启动一个代理服务
kubectl proxy
# 然后另外开一个新shell去访问一下，这个应该会返回一些内容
curl http://localhost:8001/version
```

### 查看应用的详细信息
```bash
# 首先我们可以查看所有应用
kubectl get pods
# 查看应用详细信息
kubectl describe pods
# 查看某一个应用的日志
kubectl logs kubernetes-bootcamp-65d5b99f84-4622v
# 可以执行某一个命令
kubectl exec kubernetes-bootcamp-65d5b99f84-4622v -- env
# 也可以直接进入到shell中
kubectl exec -ti kubernetes-bootcamp-65d5b99f84-4622v -- bash
```

### 暴露应用
下面我们把应用暴露出来，可以外网访问
```bash
kubectl get pods
# 获取服务信息
kubectl get services
# 会返回下面这个信息
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   32m
# 下面把这个应用暴露到公网
kubectl expose deployment/kubernetes-bootcamp --type="NodePort" --port 8080
# 然后在看一下
kubectl get services
# 返回下面的内容
NAME                  TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
kubernetes            ClusterIP   10.96.0.1      <none>        443/TCP          34m
kubernetes-bootcamp   NodePort    10.108.89.80   <none>        8080:31182/TCP   53s
# 这个时候再查看一下这个应用的信息
kubectl describe services/kubernetes-bootcamp

# 可以通过标签的形式来搜索应用
kubectl get pods -l app=kubernetes-bootcamp
# 可以手动给应用设置标签
kubectl label pods kubernetes-bootcamp-65d5b99f84-4622v version=v1
# 查看具体应用的标签信息
kubectl describe pods kubernetes-bootcamp-65d5b99f84-4622v
# 最后通过应用来获取节点信息
kubectl get pods -l version=v1

# 最后演示一下删除网络服务
kubectl delete service -l app=kubernetes-bootcamp
# 这个时候再获取一下就看不到了
kubectl get services
# 但是还是可以在容器里面访问这个应用
kubectl exec -ti kubernetes-bootcamp-65d5b99f84-4622v -- curl localhost:8080
```

### 应用缩放
应用缩放就是可以同时部署多个应用
```bash
# 查看我们部署的应用
kubectl get deployments
# 查看一下副本数
kubectl get rs
# 下面我们部署4个副本
kubectl scale deployments/kubernetes-bootcamp --replicas=4
# 然后再查看一下，这个时候就变成四个了
kubectl get deployments
# 下面查看一下这四个实例的简陋信息
kubectl get pods -o wide
# 查看一下详细信息
kubectl describe deployments/kubernetes-bootcamp
# 缩放应用，这里缩放为两个
kubectl scale deployments/kubernetes-bootcamp --replicas=2
# 这个时候再看就变成两个了
kubectl get deployments
```

### 应用更新
```bash
# 下面我们来更新应用，这里我们把镜像更新成v2的版本
kubectl set image deployments/kubernetes-bootcamp kubernetes-bootcamp=jocatalin/kubernetes-bootcamp:v2
# 然后执行下面这个命令，就可以看到一个更新的状态
kubectl get pods
# 使用下面的命令查看更新的状态
kubectl rollout status deployments/kubernetes-bootcamp
# 我们还可以回滚
kubectl rollout undo deployments/kubernetes-bootcamp
```

### 简单部署应用并给外网访问
```bash
kubectl apply -f https://k8s.io/examples/service/load-balancer-example.yaml
# 查看部署信息
kubectl get deployments hello-world
kubectl describe deployments hello-world
# 查看副本集信息
kubectl get replicasets
kubectl describe replicasets
# 创建一个server对象
kubectl expose deployment hello-world --type=LoadBalancer --name=my-service
# 获取我的服务
kubectl get services my-service
# 查看每个服务的详细信息
kubectl describe services my-service
# 查看服务的ip地址信息
kubectl get pods --output=wide

# 删除服务
kubectl delete services my-service
# 删除应用
kubectl delete deployment hello-world
```