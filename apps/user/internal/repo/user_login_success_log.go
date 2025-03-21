package repo

import (
	"context"

	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/ai-shop/user/internal/repo/query"
	"github.com/bobacgo/kit/web/r/page"
)

type UserLoginSuccessLogRepo struct {
	q *query.Query
}

func NewUserLoginSuccessLogRepo(data *Data) *UserLoginSuccessLogRepo {
	return &UserLoginSuccessLogRepo{
		q: query.Use(data.DB),
	}
}

// FindByPage 获取所有用户登录成功日志
func (r *UserLoginSuccessLogRepo) FindByPage(ctx context.Context, q *page.Query) ([]*model.UserLoginSuccessLog, int64, error) {
	return r.q.UserLoginSuccessLog.WithContext(ctx).FindByPage(q.Offset(), q.Limit())
}

// FindOneUserLoginSuccessLogById 根据ID获取用户登录成功日志
func (r *UserLoginSuccessLogRepo) FindOneUserLoginSuccessLogById(ctx context.Context, id uint64) (*model.UserLoginSuccessLog, error) {
	return r.q.UserLoginSuccessLog.WithContext(ctx).Where(r.q.UserLoginSuccessLog.ID.Eq(id)).First()
}

// InsertUserLoginSuccessLog 创建用户登录成功日志
func (r *UserLoginSuccessLogRepo) InsertUserLoginSuccessLog(ctx context.Context, log *model.UserLoginSuccessLog) error {
	return r.q.UserLoginSuccessLog.WithContext(ctx).Create(log)
}

// UpdateUserLoginSuccessLog 更新用户登录成功日志
func (r *UserLoginSuccessLogRepo) UpdateUserLoginSuccessLog(ctx context.Context, log *model.UserLoginSuccessLog) error {
	_, err := r.q.UserLoginSuccessLog.WithContext(ctx).Where(r.q.UserLoginSuccessLog.ID.Eq(log.ID)).Updates(log)
	return err
}

// DeleteUserLoginSuccessLog 删除用户登录成功日志
func (r *UserLoginSuccessLogRepo) DeleteUserLoginSuccessLog(ctx context.Context, id uint64) error {
	_, err := r.q.UserLoginSuccessLog.WithContext(ctx).Where(r.q.UserLoginSuccessLog.ID.Eq(id)).Delete()
	return err
}
