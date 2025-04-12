package config

import "github.com/bobacgo/kit/app/conf"

func Cfg() *conf.App[Service] {
	return &conf.App[Service]{
		Basic:   conf.GetBasicConf(),
		Service: conf.GetServiceConf[Service](),
	}
}

type Service struct {
	Endpoint Endpoint `mapstructure:"endpoint" yaml:"endpoint"` // 服务的grpc地址
}

type Endpoint struct {
	UserServer string `mapstructure:"userServer" yaml:"userServer"` // user服务的grpc地址
}
