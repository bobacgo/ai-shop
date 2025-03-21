package repo

import (
	"context"

	"github.com/bobacgo/ai-shop/user/internal/repo/model"
	"github.com/bobacgo/ai-shop/user/internal/repo/query"
	"github.com/bobacgo/kit/web/r/page"
)

type MerchantRepo struct {
	q *query.Query
}

func NewMerchantRepo(data *Data) *MerchantRepo {
	return &MerchantRepo{
		q: query.Use(data.DB),
	}
}

// FindByPage 获取所有商户
func (r *MerchantRepo) FindByPage(ctx context.Context, q *page.Query) ([]*model.Merchant, int64, error) {
	return r.q.Merchant.WithContext(ctx).FindByPage(q.Offset(), q.Limit())
}

// FindOneMerchantById 根据ID获取商户
func (r *MerchantRepo) FindOneMerchantById(ctx context.Context, id string) (*model.Merchant, error) {
	return r.q.Merchant.WithContext(ctx).Where(r.q.Merchant.ID.Eq(id)).First()
}

// InsertMerchant 创建商户
func (r *MerchantRepo) InsertMerchant(ctx context.Context, merchant *model.Merchant) error {
	return r.q.Merchant.WithContext(ctx).Create(merchant)
}

// UpdateMerchant 更新商户信息
func (r *MerchantRepo) UpdateMerchant(ctx context.Context, merchant *model.Merchant) error {
	_, err := r.q.Merchant.WithContext(ctx).Where(r.q.Merchant.ID.Eq(merchant.ID)).Updates(merchant)
	return err
}

// DeleteMerchant 删除商户
func (r *MerchantRepo) DeleteMerchant(ctx context.Context, id string) error {
	_, err := r.q.Merchant.WithContext(ctx).Where(r.q.Merchant.ID.Eq(id)).Delete()
	return err
}
