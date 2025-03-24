package repo

import (
	"context"
	"fmt"

	v1 "github.com/bobacgo/ai-shop/api/gen/go/user/v1"
	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/ai-shop/user/internal/repo/query"
	"github.com/bobacgo/kit/web/r/page"
	"gorm.io/gorm"
)

type UserRepo struct {
	q *query.Query
}

func NewUserRepo(data *Data) *UserRepo {
	return &UserRepo{
		q: query.Use(data.DB),
	}
}

// FindByPage 获取所有用户
func (r *UserRepo) FindByPage(ctx context.Context, q *page.Query) ([]*model.User, int64, error) {
	return r.q.User.WithContext(ctx).FindByPage(q.Offset(), q.Limit())
}

// FindOneUserById 根据ID获取用户
func (r *UserRepo) FindOneUserById(ctx context.Context, id string) (*model.User, error) {
	return r.q.User.WithContext(ctx).Where(r.q.User.ID.Eq(id)).First()
}

// InsertUser 创建用户
func (r *UserRepo) InsertUser(ctx context.Context, user *model.User) error {
	return r.q.User.WithContext(ctx).Create(user)
}

// UpdateUser 更新用户信息
func (r *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := r.q.User.WithContext(ctx).Where(r.q.User.ID.Eq(user.ID)).Updates(user)
	return err
}

// DeleteUser 删除用户
func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	_, err := r.q.User.WithContext(ctx).Where(r.q.User.ID.Eq(id)).Delete()
	return err
}

// FindOneUserByUsername 根据用户名查询用户
func (r *UserRepo) FindOneUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.q.User.WithContext(ctx).Where(r.q.User.Username.Eq(username)).First()
}

// UpdateStatus 更新用户状态
func (r *UserRepo) UpdateStatus(ctx context.Context, id string, status int32) error {
	_, err := r.q.User.WithContext(ctx).Where(r.q.User.ID.Eq(id)).Update(r.q.User.Status, status)
	return err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id string, passwd string) any {
	_, err := r.q.User.WithContext(ctx).Where(r.q.User.ID.Eq(id)).Update(r.q.User.Password, passwd)
	return err
}

// DeletedUser 注销用户
func (r *UserRepo) DeletedUser(ctx context.Context, id string) error {
	return r.q.Transaction(func(tx *query.Query) error {
		info, err := tx.User.WithContext(ctx).Where(r.q.User.ID.Eq(id)).Update(r.q.User.Status, v1.UserStatus_deleted)
		if err != nil {
			return fmt.Errorf("user status update failed: %w", err)
		} else if info.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err = tx.UserDeletionRequest.WithContext(ctx).Create(&model.UserDeletionRequest{
			UserID: id,
		}); err != nil {
			return fmt.Errorf("user deletion request create failed: %w", err)
		}
		return nil
	})
}
