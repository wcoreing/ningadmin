package main

import (
	"gorm.io/gen"
	"wnadmin/app/global"
	"wnadmin/app/provider/conf"
	"wnadmin/app/provider/db"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	conf.Init()
	db.Init()
	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(global.DB.QPGameUserDB) // reuse your gorm db
	AccountsPlatform := g.GenerateModel("Accounts_platform")
	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(AccountsPlatform)

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	g.ApplyInterface(func(Querier) {}, AccountsPlatform)

	// Generate the code
	g.Execute()
}
