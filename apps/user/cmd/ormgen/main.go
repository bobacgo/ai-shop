package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ai_shop_user"))
	if err != nil {
		panic(err)
	}
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../user/internal/repo/query",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		FieldCoverable:    false,
	})
	g.UseDB(db)
	g.ApplyBasic(
		g.GenerateModel("merchants"),
		g.GenerateModel("tags"),
		g.GenerateModel("user_addresses"),
		g.GenerateModel("user_points"),
		g.GenerateModel("user_profiles"),
		g.GenerateModel("users"),
	)
	g.Execute()
}
