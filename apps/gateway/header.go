package main

import (
	"context"
	"net/http"

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
	return metadata.New(md)
}
