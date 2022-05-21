## 搜索引擎
参考： https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html
```bash
sudo docker pull docker.elastic.co/elasticsearch/elasticsearch:8.2.0
这个服务需要 9200 和9300两个端口，指定下面两个环境变量
# 环境变量
discovery.type=single-node
ES_JAVA_OPTS="-Xms200m -Xmx200m"
# 配置文件地址
/usr/share/elasticsearch/config
/usr/share/elasticsearch/data
/usr/share/elasticsearch/plugins
# es的管理工具 kibana，端口为5601
sudo docker pull docker.elastic.co/kibana/kibana:8.2.0
# 指定环境变量(对应9200)
ELASTICSEARCH_HOSTS = http://192.168.1.50:30003
```

如果想映射es的配置文件需要自己不映射，然后把配置文件给拷贝出来，配置文件如下，配置文件在 `/usr/share/elasticsearch/config/elasticsearch.yml` 中
```
cluster.name: "docker-cluster"
network.host: 0.0.0.0
# 关闭密码
xpack.security.enabled: false
```
