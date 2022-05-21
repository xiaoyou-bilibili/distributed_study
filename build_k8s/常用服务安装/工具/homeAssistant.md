
```bash
sudo docker pull homeassistant/home-assistant
```

注意，如果访问报错400bad request，可以自己看一下日志，查看一下错误信息，然后在`config\configuration.yaml` 里面配置一下
```yaml
http:
  use_x_forwarded_for: true
  trusted_proxies:
    - 100.0.0.0/8  # Add the IP address of the proxy server
```