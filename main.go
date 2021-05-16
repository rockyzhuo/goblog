package main

import (
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
