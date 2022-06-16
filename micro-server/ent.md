这里我们使用facebook开源的数据库操作框架

文档地址：https://entgo.io/zh/docs/getting-started

## 简单使用
```bash
# 安装工具
go get entgo.io/ent/cmd/ent
# 初始化我们的实体
ent init User
# 也可以用下面这个命令
go run entgo.io/ent/cmd/ent init User

# 初始化好的实体会在 project/ent/schema/ 目录下
# 我们可以自己去配置字段



```