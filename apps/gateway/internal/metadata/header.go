package metadata

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

const (
	Language = "language"
	Platform = "platform"
)

func HeaderToMD(ctx context.Context, req *http.Request) metadata.MD {
	// 处理请求头，将HTTP请求头转换为gRPC metadata
	md := make(map[string]string)
	if auth := req.Header.Get(Language); auth != "" {
		md[Language] = auth
	}
	if auth := req.Header.Get(Platform); auth != "" {
		md[Platform] = auth
	}

	if method, ok := runtime.RPCMethod(ctx); ok {
		md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
	}
	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		md["pattern"] = pattern // /v1/example/login
	}
	return metadata.New(md)
}
