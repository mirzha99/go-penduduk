package main

import (
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/router"
)

// @title 	Tag Service API
// @version	1.0
// @description A Tag service API in Go using Gin framework

// @host 	localhost:8888
// @BasePath /api
func main() {
	//call .env
	config.Loadenv()
	//call database gorm driver mysql
	config.ConnectionDB()
	//run gin framework
	router.Router()
}
