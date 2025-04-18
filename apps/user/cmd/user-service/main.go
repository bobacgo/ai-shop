package main

import (
	"flag"
	"log"

	"github.com/bobacgo/ai-shop/user/internal/config"
	"github.com/bobacgo/ai-shop/user/internal/server"
	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/app/conf"
)

var filepath = flag.String("config", "./config.yaml", "config file path")

func init() {
	flag.String("name", "user-service", "service name")
	flag.String("env", "dev", "run config context")
	flag.String("logger.level", "info", "logger level")
	flag.Int("port", 8080, "http port 8080, rpc port 9080")
	conf.BindPFlags()
}

func main() {
	newApp := app.New[config.Service](*filepath,
		app.WithTracerServer(),
		app.WithMustDB(),
		app.WithMustRedis(),
		app.WithGrpcServer(server.GrpcRegisterServer),
	)
	if err := newApp.Run(); err != nil {
		log.Panic(err.Error())
	}
}
