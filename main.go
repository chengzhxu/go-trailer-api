package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haxqer/gintools/logging"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "go-trailer-api/docs"
	"go-trailer-api/pkg/gredis"
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/setting"
	"go-trailer-api/pkg/tool"
	"go-trailer-api/pkg/validator"
	"go-trailer-api/routers"
	"log"
	"net/http"
)

const (
	projectName = "TrailerApi"
)

func init() {
	logsPath := tool.GetLogsPath() + projectName
	setting.Setup()
	logging.Setup(logsPath)
	validator.Setup()
	gredis.Setup()
	model.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
