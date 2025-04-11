package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// 代理配置结构体 Proxy configuration struct
type ProxyConfig struct {
	PathPrefix  string // 请求路径前缀 Path prefix for requests
	TargetURL   string // 目标服务器地址 Target server address
	StripPrefix bool   // 是否移除前缀 Whether to remove prefix
}

// 创建反向代理处理器 Create reverse proxy handler
func createReverseProxy(config ProxyConfig) (http.Handler, error) {
	targetURL, err := url.Parse(config.TargetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 修改代理的导演功能 Modify the proxy director function
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// 如果需要移除前缀 If prefix should be removed
		if config.StripPrefix {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, config.PathPrefix)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
		}

		// 设置代理头 Set proxy headers
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", targetURL.Host)
	}

	return proxy, nil
}

func main() {
	// 定义代理配置 Define proxy configurations
	proxyConfigs := []ProxyConfig{
		{
			PathPrefix:  "/api/",
			TargetURL:   "http://localhost:8080",
			StripPrefix: true,
		},
	}

	// 设置代理路由 Set up proxy routes
	for _, config := range proxyConfigs {
		proxyHandler, err := createReverseProxy(config)
		if err != nil {
			log.Fatalf("创建代理失败 Failed to create proxy for %s: %v", config.PathPrefix, err)
		}

		// 注册代理处理器 Register proxy handler
		http.HandleFunc(config.PathPrefix, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("代理请求 Proxying request: %s -> %s%s", r.URL.Path, config.TargetURL, r.URL.Path)
			proxyHandler.ServeHTTP(w, r)
		})
	}

	// 提供静态文件服务 Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("服务已启动 Server started at http://localhost:80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("服务启动失败 Server failed to start: %v", err)
	}
}
