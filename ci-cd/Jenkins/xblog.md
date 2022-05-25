## 博客系统自动化构建整理

### jenkins调用远程docker
> 因为我自己的jenkins是部署在容器里的，不能直接调用docker，所以需要连接远程docker

需要的插件：
- Publish over SSH (项目远程部署)
- Hudson SCP publisher （拷贝文件）
- Publish Over SSH （执行命令）