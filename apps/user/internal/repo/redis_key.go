package repo

import (
	"fmt"

	"github.com/bobacgo/ai-shop/user/internal/config"
)

// redis key
// 模块:主题:xxx:数据结构类型

func captchaKey(id string) string {
	return fmt.Sprintf("%s:captcha:%s:string", config.Cfg().Name, id)
}
