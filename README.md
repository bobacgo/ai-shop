# gogo-services

![](./image.png)

### 项目组件依赖

- [PolarisMesh](https://polarismesh.cn) DNS 注册中心、配置中心
- [Gin](https://gin-gonic.com) Web Framework
- [gRPC](https://grpc.io) RPC Framework
- [GORM](https://gorm.io) ORM library，操作关系型数据库
- [go-redis](https://redis.uptrace.dev) Golang Redis client for Redis Server and Redis Cluster

```text
gogo-services
    ├─admin-service -- 平台运营服务、用户中心
    ├─framework    -- 基础库
    ├─devops-conf   -- 部署相关
    ├─gateway       -- 后端网关
    ├─gen           -- 代码生成
    └─main-service  -- 主服务
```

### create project

> 1. 守护进程的项目以 -service 结尾
> 2. 每个子项目都要有 README.md 文件

```shell
# first create
go work init
```

```shell
# add project
cd ./xxx-service
go mod init github.com/gogoclouds/gogo-services/xxx-service
cd ../
go work use ./xxx-service
```

### Develop Env

部署文档参考 [devops-conf](./devops-conf) 模块

polaris、mysql、redis

#### local hosts config file

[hosts](./devops-conf/hosts)

### Reference

- [https://microservices.io](https://microservices.io)