## nginx缓存优化

https://segmentfault.com/a/1190000019179879
https://www.cnblogs.com/itzgr/p/13321980.html#_label0_1

```nginx
# 设置缓存路径
proxy_cache_path /tmp/nginx/cache levels=1:2 keys_zone=my_cache:10m max_size=10g inactive=60m use_temp_path=off;


#gzip  on;
server {
    listen 80;
    server_name  index.xiaoyou.host;# 服务器地址或绑定域名

    location /data {
    #proxy_pass http://index-data-0.index-data:8000;
    
    # 打开缓存
    proxy_cache my_cache;
    # 只缓存get请求
    proxy_cache_methods GET;
    # 缓存状态码（状态码缓存12H）
    proxy_cache_valid 200 304 301 302 12h;
    # 设置是本地缓存30天
    expires 12h;
    # 支持分片（支持audio标签拖动）
    add_header Accept-Ranges bytes;
    
    proxy_pass http://goland-0.goland.xiaoyou-tool.svc.cluster.local:8001;
    
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "Upgrade";
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header Host $host;
    }
```