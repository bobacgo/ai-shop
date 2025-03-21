package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type CaptchaRepo struct {
	rdb redis.UniversalClient
}

func NewCaptchaRepo(data *Data) *CaptchaRepo {
	return &CaptchaRepo{
		rdb: data.Rds,
	}
}

func (repo *CaptchaRepo) Set(id string, value string) error {
	err := repo.rdb.Set(context.Background(), captchaKey(id), value, time.Minute).Err()
	return err
}

func (repo *CaptchaRepo) Get(id string, clear bool) string {
	key := captchaKey(id)
	ctx := context.Background()
	result, err := repo.rdb.Get(ctx, key).Result()
	if err != nil {
		slog.Error("get captcha error", "key", key, "err", err)
		return ""
	}
	if clear {
		if err = repo.rdb.Del(ctx, key).Err(); err != nil {
			slog.Error("del captcha error", "key", key, "err", err)
			return ""
		}
	}
	return result
}

func (repo *CaptchaRepo) Verify(id, answer string, clear bool) bool {
	captchaCode := repo.Get(id, clear)
	return captchaCode == answer
}
