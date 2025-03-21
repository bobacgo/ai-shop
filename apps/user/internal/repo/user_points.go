package repo

import (
	"context"

	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/ai-shop/user/internal/repo/query"
	"github.com/bobacgo/kit/web/r/page"
)

type UserPointRepo struct {
	q *query.Query
}

func NewUserPointRepo(data *Data) *UserPointRepo {
	return &UserPointRepo{
		q: query.Use(data.DB),
	}
}

// FindByPage 获取所有用户积分记录
func (r *UserPointRepo) FindByPage(ctx context.Context, q *page.Query) ([]*model.UserPoint, int64, error) {
	return r.q.UserPoint.WithContext(ctx).FindByPage(q.Offset(), q.Limit())
}

// FindOneUserPointById 根据ID获取用户积分记录
func (r *UserPointRepo) FindOneUserPointById(ctx context.Context, id string) (*model.UserPoint, error) {
	return r.q.UserPoint.WithContext(ctx).Where(r.q.UserPoint.ID.Eq(id)).First()
}

// InsertUserPoint 创建用户积分记录
func (r *UserPointRepo) InsertUserPoint(ctx context.Context, userPoint *model.UserPoint) error {
	return r.q.UserPoint.WithContext(ctx).Create(userPoint)
}

// UpdateUserPoint 更新用户积分记录
func (r *UserPointRepo) UpdateUserPoint(ctx context.Context, userPoint *model.UserPoint) error {
	_, err := r.q.UserPoint.WithContext(ctx).Where(r.q.UserPoint.ID.Eq(userPoint.ID)).Updates(userPoint)
	return err
}

// DeleteUserPoint 删除用户积分记录
func (r *UserPointRepo) DeleteUserPoint(ctx context.Context, id string) error {
	_, err := r.q.UserPoint.WithContext(ctx).Where(r.q.UserPoint.ID.Eq(id)).Delete()
	return err
}
