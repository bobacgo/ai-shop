package register

import (
	"context"
	"log"

	v1 "github.com/bobacgo/ai-shop/api/gen/go/user/v1"
	"github.com/bobacgo/ai-shop/gatway/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Handler(mux *runtime.ServeMux) {
	// 设置grpc连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}

	endpoint := config.Cfg().Service.Endpoint

	ctx := context.Background()
	if err := v1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint.UserServer, opts); err != nil {
		log.Fatalln(err)
	}
	if err := v1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint.UserServer, opts); err != nil {
		log.Fatalln(err)
	}
	if err := v1.RegisterMerchantServiceHandlerFromEndpoint(ctx, mux, endpoint.UserServer, opts); err != nil {
		log.Fatalln(err)
	}
}
