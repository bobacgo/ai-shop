package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/bobacgo/ai-shop/api/gen/go/user/v1"
	"github.com/bobacgo/ai-shop/gatway/middleware"
)

const (
	grpcServerEndpoint = "user-service.default.svc.cluster.local" // user服务的gRPC地址
	gatewayPort        = 8080
)

func main() {
	// Initialize tracing and handle the tracer provider shutdown
	stopTracing := initTracing()
	defer stopTracing()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建grpc-gateway的mux
	mux := runtime.NewServeMux(
		runtime.WithMiddlewares(middleware.Middlewares...),
		runtime.WithMetadata(HeaderToMD),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	// 设置grpc连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}

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
	slog.Info("Server listening on ...", "port", gatewayPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", gatewayPort), mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
