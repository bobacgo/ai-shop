package repo

import (
	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/app/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Data struct {
	Cache cache.Cache
	Rds   redis.UniversalClient
	DB    *gorm.DB
}

func NewData(app *app.AppOptions) *Data {
	return &Data{
		Cache: app.LocalCache(),
		Rds:   app.Redis(),
		DB:    app.DB().Default(),
	}
}
