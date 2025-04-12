package main

import (
	"flag"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/bobacgo/ai-shop/gatway/internal/config"
	"github.com/bobacgo/ai-shop/gatway/internal/metadata"
	"github.com/bobacgo/ai-shop/gatway/internal/middleware"
	"github.com/bobacgo/ai-shop/gatway/internal/register"
	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/app/conf"
)

var filepath = flag.String("config", "./config.yaml", "config file path")

func init() {
	flag.String("name", "user-service", "service name")
	flag.String("env", "dev", "run config context")
	flag.String("logger.level", "info", "logger level")
	flag.Int("port", 8080, "http port 8080")
	conf.BindPFlags()
}

func main() {
	// Initialize tracing and handle the tracer provider shutdown
	// stopTracing := trace.InitTracing()
	// defer stopTracing()

	newApp := app.New[config.Service](*filepath,
		app.WithMustRedis(),
		app.WithGatewayServer(register.Handler,
			runtime.WithMiddlewares(middleware.Middlewares...),
			runtime.WithMetadata(metadata.HeaderToMD),
			runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		),
	)
	if err := newApp.Run(); err != nil {
		log.Panic(err.Error())
	}
}
