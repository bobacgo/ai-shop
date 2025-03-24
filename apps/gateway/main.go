package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	v1 "github.com/bobacgo/ai-shop/api/pb/user/v1"
)

const (
	grpcServerEndpoint = "localhost:9080" // user服务的gRPC地址
	gatewayPort        = 8080
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建grpc-gateway的mux
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			// 处理请求头，将HTTP请求头转换为gRPC metadata
			md := make(map[string]string)
			if auth := req.Header.Get("Authorization"); auth != "" {
				md["authorization"] = auth
			}
			return metadata.New(md)
		}),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	// 设置grpc连接选项
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 注册用户服务的HTTP处理器
	if err := v1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts); err != nil {
		log.Fatalf("Failed to register auth service handler: %v", err)
	}
	if err := v1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts); err != nil {
		log.Fatalf("Failed to register user service handler: %v", err)
	}
	if err := v1.RegisterMerchantServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts); err != nil {
		log.Fatalf("Failed to register merchant service handler: %v", err)
	}
	// 启动HTTP服务器
	log.Printf("Server listening on port %d...", gatewayPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", gatewayPort), mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}