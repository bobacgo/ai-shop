package config

import "github.com/bobacgo/kit/app/conf"

func Cfg() Service {
	return conf.GetServiceConf[Service]()
}

type Service struct {
	ErrAttemptLimit int `mapstructure:"errAttemptLimit" yaml:"errAttemptLimit" mask:""`
}