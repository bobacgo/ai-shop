package middleware

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

var Middlewares = []runtime.Middleware{
	LoggingMiddleware,
	// AuthMiddleware,
}
