package repo

import (
	"context"

	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/ai-shop/user/internal/repo/query"
	"github.com/bobacgo/kit/web/r/page"
)

type UserAddressRepo struct {
	q *query.Query
}

func NewUserAddressRepo(data *Data) *UserAddressRepo {
	return &UserAddressRepo{
		q: query.Use(data.DB),
	}
}

// FindByPage 获取所有用户地址
func (r *UserAddressRepo) FindByPage(ctx context.Context, q *page.Query) ([]*model.UserAddress, int64, error) {
	return r.q.UserAddress.WithContext(ctx).FindByPage(q.Offset(), q.Limit())
}

// FindOneUserAddressById 根据ID获取用户地址
func (r *UserAddressRepo) FindOneUserAddressById(ctx context.Context, id string) (*model.UserAddress, error) {
	return r.q.UserAddress.WithContext(ctx).Where(r.q.UserAddress.ID.Eq(id)).First()
}

// InsertUserAddress 创建用户地址
func (r *UserAddressRepo) InsertUserAddress(ctx context.Context, userAddress *model.UserAddress) error {
	return r.q.UserAddress.WithContext(ctx).Create(userAddress)
}

// UpdateUserAddress 更新用户地址信息
func (r *UserAddressRepo) UpdateUserAddress(ctx context.Context, userAddress *model.UserAddress) error {
	_, err := r.q.UserAddress.WithContext(ctx).Where(r.q.UserAddress.ID.Eq(userAddress.ID)).Updates(userAddress)
	return err
}

// DeleteUserAddress 删除用户地址
func (r *UserAddressRepo) DeleteUserAddress(ctx context.Context, id string) error {
	_, err := r.q.UserAddress.WithContext(ctx).Where(r.q.UserAddress.ID.Eq(id)).Delete()
	return err
}
