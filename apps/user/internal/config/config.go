package config

import "github.com/bobacgo/kit/app/conf"

func Cfg() *conf.App[Service] {
	return &conf.App[Service]{
		Basic:   conf.GetBasicConf(),
		Service: conf.GetServiceConf[Service](),
	}
}

type Service struct {
	ErrAttemptLimit int `mapstructure:"errAttemptLimit" yaml:"errAttemptLimit" mask:""`
}
