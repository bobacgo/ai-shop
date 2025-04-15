# ai-shop

![img](./gzh.png)

```
.
├── LICENSE
├── README.md
├── TODO.md
├── api
│   ├── Makefile
│   ├── buf.gen.yaml
│   ├── buf.lock
│   ├── buf.yaml
│   ├── cmd
│   │   ├── errgen
│   │   │   └── main.go
│   │   └── protogen
│   │       └── main.go
│   ├── gen
│   │   ├── go
│   │   │   ├── cart
│   │   │   │   └── v1
│   │   │   │       └── errs
│   │   │   │           ├── cart_error.pb.go
│   │   │   │           └── error_message.gen.go
│   │   │   └── user
│   │   │       └── v1
│   │   │           ├── auth.pb.go
│   │   │           ├── auth.pb.gw.go
│   │   │           ├── auth_grpc.pb.go
│   │   │           ├── errs
│   │   │           │   ├── error_message.gen.go
│   │   │           │   └── user_error.pb.go
│   │   │           ├── user.pb.go
│   │   │           ├── user.pb.gw.go
│   │   │           └── user_grpc.pb.go
│   │   └── openapi
│   │       └── v1
│   │           ├── auth.swagger.json
│   │           ├── cart_error.swagger.json
│   │           ├── user.swagger.json
│   │           └── user_error.swagger.json
│   ├── go.mod
│   ├── go.sum
│   └── proto
│       ├── cart
│       │   └── v1
│       │       └── cart_error.proto
│       ├── order
│       ├── payment
│       ├── product
│       ├── third_party
│       │   └── google
│       │       └── api
│       │           ├── annotations.proto
│       │           └── http.proto
│       └── user
│           └── v1
│               ├── auth.proto
│               ├── user.proto
│               └── user_error.proto
├── apps
│   ├── cart
│   │   └── go.mod
│   ├── frontend
│   │   ├── go.mod
│   │   ├── nginx.go
│   │   └── static
│   │       ├── css
│   │       │   └── common.css
│   │       ├── image
│   │       │   └── gzh.png
│   │       ├── login.html
│   │       └── register.html
│   ├── gateway
│   │   ├── cmd
│   │   │   └── gateway
│   │   │       └── main.go
│   │   ├── config.yaml
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── internal
│   │       ├── config
│   │       │   └── config.go
│   │       ├── metadata
│   │       │   └── header.go
│   │       ├── middleware
│   │       │   ├── auth.go
│   │       │   ├── ctx.go
│   │       │   ├── log.go
│   │       │   └── middeware.go
│   │       ├── register
│   │          └── register.go
│   ├── inventory
│   │   └── go.mod
│   ├── order
│   │   └── go.mod
│   ├── payment
│   │   └── go.mod
│   ├── product
│   │   └── go.mod
│   └── user
│       ├── Makefile
│       ├── cmd
│       │   ├── ormgen
│       │   │   └── main.go
│       │   └── user-service
│       │       └── main.go
│       ├── config.yaml
│       ├── deploy
│       │   ├── configmap.yaml
│       │   └── k8s.yaml
│       ├── go.mod
│       ├── go.sum
│       └── internal
│           ├── config
│           │   └── config.go
│           ├── repo
│           │   ├── captcha.go
│           │   ├── data.go
│           │   ├── merchants.go
│           │   ├── model
│           │   │   ├── merchants.gen.go
│           │   │   ├── tags.gen.go
│           │   │   ├── user_addresses.gen.go
│           │   │   ├── user_deletion_requests.gen.go
│           │   │   ├── user_login_success_log.gen.go
│           │   │   ├── user_points.gen.go
│           │   │   ├── user_profiles.gen.go
│           │   │   └── users.gen.go
│           │   ├── query
│           │   │   ├── gen.go
│           │   │   ├── merchants.gen.go
│           │   │   ├── tags.gen.go
│           │   │   ├── user_addresses.gen.go
│           │   │   ├── user_deletion_requests.gen.go
│           │   │   ├── user_login_success_log.gen.go
│           │   │   ├── user_points.gen.go
│           │   │   ├── user_profiles.gen.go
│           │   │   └── users.gen.go
│           │   ├── redis_key.go
│           │   ├── repo.go
│           │   ├── user.go
│           │   ├── user_addresses.go
│           │   ├── user_deletion_requests.go
│           │   ├── user_login_success_log.go
│           │   ├── user_points.go
│           │   └── user_profiles.go
│           ├── server
│           │   └── grpc.go
│           └── service
│               ├── auth.go
│               └── user.go
├── deploy
│   ├── Makefile
│   ├── README.md
│   ├── ai-shop
│   │   └── v1.0.0
│   │       ├── common.yaml
│   │       ├── db.sql
│   │       ├── gateway.yaml
│   │       ├── mysql.yaml
│   │       └── redis.yaml
│   ├── config.yaml
│   ├── development_tool.md
│   ├── docker
│   │   ├── mysql
│   │   │   ├── config.md
│   │   │   ├── my.cnf
│   │   │   ├── mysql.sh
│   │   │   └── tree.png
│   │   ├── polaris.sh
│   │   └── redis
│   │       ├── redis.conf
│   │       └── redis.sh
│   ├── etcd
│   ├── hosts
│   └── sql.dbcnb
├── go.work
├── go.work.sum
└── gzh.png
```

### Reference

- [https://microservices.io](https://microservices.io)
