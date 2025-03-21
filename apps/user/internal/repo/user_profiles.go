package repo

import (
	"context"

	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/ai-shop/user/internal/repo/query"
	"github.com/bobacgo/kit/web/r/page"
)

type UserProfileRepo struct {
	q *query.Query
}

func NewUserProfileRepo(data *Data) *UserProfileRepo {
	return &UserProfileRepo{
		q: query.Use(data.DB),
	}
}

// FindByPage 获取所有用户详情
func (r *UserProfileRepo) FindByPage(ctx context.Context, q *page.Query) ([]*model.UserProfile, int64, error) {
	return r.q.UserProfile.WithContext(ctx).FindByPage(q.Offset(), q.Limit())
}

// FindOneUserProfileById 根据ID获取用户详情
func (r *UserProfileRepo) FindOneUserProfileById(ctx context.Context, id string) (*model.UserProfile, error) {
	return r.q.UserProfile.WithContext(ctx).Where(r.q.UserProfile.UserID.Eq(id)).First()
}

// InsertUserProfile 创建用户详情
func (r *UserProfileRepo) InsertUserProfile(ctx context.Context, userProfile *model.UserProfile) error {
	return r.q.UserProfile.WithContext(ctx).Create(userProfile)
}

// UpdateUserProfile 更新用户详情信息
func (r *UserProfileRepo) UpdateUserProfile(ctx context.Context, userProfile *model.UserProfile) error {
	_, err := r.q.UserProfile.WithContext(ctx).Where(r.q.UserProfile.UserID.Eq(userProfile.UserID)).Updates(userProfile)
	return err
}

// DeleteUserProfile 删除用户详情
func (r *UserProfileRepo) DeleteUserProfile(ctx context.Context, id string) error {
	_, err := r.q.UserProfile.WithContext(ctx).Where(r.q.UserProfile.UserID.Eq(id)).Delete()
	return err
}
