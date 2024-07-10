package main

import (
	"embed"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	"net/http"
)

//go:embed resources/views/articles/*
//go:embed resources/views/auth/*
//go:embed resources/views/categories/*
//go:embed resources/views/layouts/*
var tplFS embed.FS

func init() {
	config.Initialize()
}
func main() {
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	bootstrap.SetupTemplate(tplFS)

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
