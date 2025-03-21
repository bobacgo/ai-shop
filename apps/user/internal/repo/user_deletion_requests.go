package repo

import (
	"context"

	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/ai-shop/user/internal/repo/query"
	"github.com/bobacgo/kit/web/r/page"
)

type UserDeletionRequestRepo struct {
	q *query.Query
}

func NewUserDeletionRequestRepo(data *Data) *UserDeletionRequestRepo {
	return &UserDeletionRequestRepo{
		q: query.Use(data.DB),
	}
}

// FindByPage 获取所有用户注销请求记录
func (r *UserDeletionRequestRepo) FindByPage(ctx context.Context, q *page.Query) ([]*model.UserDeletionRequest, int64, error) {
	return r.q.UserDeletionRequest.WithContext(ctx).FindByPage(q.Offset(), q.Limit())
}

// FindOneUserDeletionRequestById 根据ID获取用户注销请求记录
func (r *UserDeletionRequestRepo) FindOneUserDeletionRequestById(ctx context.Context, id string) (*model.UserDeletionRequest, error) {
	return r.q.UserDeletionRequest.WithContext(ctx).Where(r.q.UserDeletionRequest.ID.Eq(id)).First()
}

// InsertUserDeletionRequest 创建用户注销请求记录
func (r *UserDeletionRequestRepo) InsertUserDeletionRequest(ctx context.Context, request *model.UserDeletionRequest) error {
	return r.q.UserDeletionRequest.WithContext(ctx).Create(request)
}

// UpdateUserDeletionRequest 更新用户注销请求记录
func (r *UserDeletionRequestRepo) UpdateUserDeletionRequest(ctx context.Context, request *model.UserDeletionRequest) error {
	_, err := r.q.UserDeletionRequest.WithContext(ctx).Where(r.q.UserDeletionRequest.ID.Eq(request.ID)).Updates(request)
	return err
}

// DeleteUserDeletionRequest 删除用户注销请求记录
func (r *UserDeletionRequestRepo) DeleteUserDeletionRequest(ctx context.Context, id string) error {
	_, err := r.q.UserDeletionRequest.WithContext(ctx).Where(r.q.UserDeletionRequest.ID.Eq(id)).Delete()
	return err
}
