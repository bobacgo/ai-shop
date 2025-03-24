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
	repoAll := repo.New(
		repo.NewCaptchaRepo(data),
		repo.NewUserRepo(data),
		repo.NewUserDeletionRequestRepo(data),
		repo.NewUserLoginSuccessLogRepo(data),
	)

	// register
	v1.RegisterUserServiceServer(srv, service.NewUserService(repoAll.User))
	v1.RegisterAuthServiceServer(srv, service.NewAuthService(data.Rds, repoAll))
	// v1.RegisterMerchantServiceServer(srv, &v1.Server{})
}