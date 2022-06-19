## 镜像打包
```bash
# 拉取镜像
sudo docker pull jupyter/scipy-notebook:hub-2.3.1
# 安装c++扩展
conda install xeus-cling -c conda-forge
# 执行下面这个命令就可以了，前面是容器名字，后面是镜像的名称
sudo docker commit cc726345ce3c jupyter/scipy-notebook:hub-2.3.1-cpp
``
