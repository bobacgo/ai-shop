# Protobuf 配置
PROTO_SRC = proto  # proto 目录
PB_OUT = pb        # 生成 pb.go 代码的目录
PROTOC = protoc    # protoc 编译器
GO_OUT = --go_out=$(PB_OUT) --go-grpc_out=$(PB_OUT) --proto_path=$(PROTO_SRC) --proto_path=$(PROTO_SRC)/third_party

# 递归查找所有 .proto 文件
PROTO_FILES := $(shell find $(PROTO_SRC) -name "*.proto")

# init env
# 使用方式: make init
# 需要添加其他 tool, 请执行 go get -tool github.com/xxx/xxx@latest
.PHONY: init
init:
	go get -u google.golang.org/grpc@latest
	go get -u google.golang.org/genproto/googleapis/api
	go mod tidy
	go install ./...


# 默认构建所有 proto
# 使用方式: make proto
.PHONY: proto
proto:
	@echo "Generating all gRPC code..."
	$(PROTOC) $(GO_OUT) $(PROTO_FILES)
	@echo "All gRPC code generated successfully!"

# 生成指定模块的proto代码和错误映射文件
# 使用方式: make <module_name>
.PHONY: user cart
user cart:
	@echo "Generating $@ proto and error mapping..."
	go run cmd/protogen/main.go $@
	go run cmd/errgen/main.go $@
	@echo "$@ proto and error mapping generated successfully!"

# 清理 pb 目录
.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	find $(PB_OUT) -name "*.go" -type f -delete
	@echo "Clean completed!"