package middleware

import (
	"context"

	v1 "github.com/bobacgo/ai-shop/api/gen/go/user/v1"
)

type userContextKey struct{} // 私有 key 类型

func WithUser(ctx context.Context, user *v1.UserTokenInfo) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}

func GetUser(ctx context.Context) *v1.UserTokenInfo {
	user, _ := ctx.Value(userContextKey{}).(*v1.UserTokenInfo)
	return user
}
