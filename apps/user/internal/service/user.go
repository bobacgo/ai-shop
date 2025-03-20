package service

import (
	"context"

	v1 "github.com/bobacgo/ai-shop/api/pb/user/v1"
	"github.com/bobacgo/ai-shop/api/pb/user/v1/errs"
	"github.com/bobacgo/ai-shop/user/internal/repo"
	"github.com/bobacgo/kit/app/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	v1.UnimplementedUserServiceServer
	repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u UserService) GetUserById(ctx context.Context, request *v1.GetUserRequest) (*v1.UserResponse, error) {
	return nil, errs.Status(ctx, errs.Err_UserNotFound)
}

func (u UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.UserResponse, error) {
	if err := validator.StructCtx(ctx, request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return nil, errs.Status(ctx, errs.Err_UserNotFound)
}

func (u UserService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u UserService) DeleteUser(ctx context.Context, request *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	// TODO implement me
	return &emptypb.Empty{}, nil
}
