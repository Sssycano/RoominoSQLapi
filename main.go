package main

import (
	"roomino/dao"
	"roomino/routes"
)

func main() {

	dao.MySQLInit()
	r := routes.NewRouter()
	_ = r.Run(":3000")
}
