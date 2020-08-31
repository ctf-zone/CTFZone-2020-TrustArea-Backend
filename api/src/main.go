package main

import (
	"api/app"
	"api/app/db"
	"api/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	config.ConfigInit()
	db.DBInit()
	db.RedisInit()
	r := app.GenerateRoutes()
	addr := fmt.Sprintf("%s:%d", config.Config.App.Host, config.Config.App.Port)
	log.Infof("Started API server: http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}