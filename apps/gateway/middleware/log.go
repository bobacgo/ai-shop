package middleware

import (
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// 日志中间件
func LoggingMiddleware(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next(w, r, pathParams) // 调用下一个 HandlerFunc
	}
}

// TODO 失败响应日志
