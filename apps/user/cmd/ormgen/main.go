package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/ai_shop_user"))
	if err != nil {
		panic(err)
	}
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../user/internal/repo/query", // 生成的查询代码存放在 ./query 目录
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,  // NULL 字段会变成 *string 而不是 string
		FieldSignable:     true,  // unsigned int 识别并匹配到正确的 Go 类型
		FieldWithIndexTag: true,  // 结构体字段添加 gorm:"index"，方便查询优化
		FieldWithTypeTag:  true,  // 结构体字段添加 gorm:"type:varchar(255)"，保留数据库类型信息
		FieldCoverable:    false, // 默认值字段生成 *int，防止 0 被错误解析
	})
	g.UseDB(db)
	g.ApplyBasic(
		g.GenerateModel("merchants"),
		g.GenerateModel("tags"),
		g.GenerateModel("user_addresses"),
		g.GenerateModel("user_points"),
		g.GenerateModel("user_profiles"),
		g.GenerateModel("users"),
		g.GenerateModel("user_deletion_requests"),
		g.GenerateModel("user_login_success_log"),
	)
	g.Execute()
}
