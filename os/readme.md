## 常用服务器操作系统


## 常见问题
### 设置root密码
```bash
sudo passwd root
```
### ssh允许root账号

`sudo vim /etc/ssh/sshd_config`

调整PermitRootLogin参数值为yes，如下图：
![](../images/2022-05-15-20-24-56.png)

```bash
PermitRootLogin yes
```

重启ssh服务 `sudo  systemctl  restart  ssh`