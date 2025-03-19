package server

import (
	v1 "github.com/bobacgo/ai-shop/api/pb/user/v1"
	"github.com/bobacgo/ai-shop/user/internal/repo"
	"github.com/bobacgo/ai-shop/user/internal/service"
	"github.com/bobacgo/kit/app"
	"google.golang.org/grpc"
)

func GrpcRegisterServer(srv *grpc.Server, comps *app.AppOptions) {
	data := repo.NewData(comps)

	// repo
	userRepo := repo.NewUserRepo(data)

	// register
	v1.RegisterUserServiceServer(srv, service.NewUserService(userRepo))
	// v1.RegisterMerchantServiceServer(srv, &v1.Server{})
}