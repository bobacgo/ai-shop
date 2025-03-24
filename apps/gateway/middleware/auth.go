package middleware

import (
	"net/http"
	"slices"

	v1 "github.com/bobacgo/ai-shop/api/gen/go/user/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var whiteList = []string{
	"v1/auth/login",
	"v1/auth/register",
	"v1/auth/send_verification_code",
	"v1/auth/reset_password",
	"v1/auth/refresh_token",
}

// 认证中间件
func AuthMiddleware(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		if slices.Contains(whiteList, r.URL.Path) {
			next(w, r, pathParams)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO 解析，并验证 token
		if token != "Bearer valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userInfo := &v1.UserTokenInfo{
			UserId:   "10000001",
			Username: "admin",
			Role:     v1.Role_admin,
		}

		// 将用户信息存入 context
		ctx := WithUser(r.Context(), userInfo)
		next(w, r.WithContext(ctx), pathParams) // 调用下一个 HandlerFunc
	}
}
