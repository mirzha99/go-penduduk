package main

import (
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/router"
)

// @title Go Penduduk
// @version	1.0
// @description A Go Penduduk in Go using Gin framework

// @host 	localhost:3131
// @BasePath /
func main() {
	//call .env
	config.Loadenv()
	//call database gorm driver mysql
	config.ConnectionDB()
	//run gin framework
	router.Router()
}
