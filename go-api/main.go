package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/peashoot/peapi/avatar"
	"github.com/peashoot/peapi/cnarea"
	"github.com/peashoot/peapi/config"
	_ "github.com/peashoot/peapi/docs"
	"github.com/peashoot/peapi/webhook"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	flag.StringVar(&confPath, "c", "conf.json", "path of configuration file")
}

var (
	confPath string // 配置文件路径
)

func main() {
	config.ReadConfig(confPath)
	r := gin.Default()
	r.POST("/cnarea", cnarea.Handle)
	r.POST("/webhook", webhook.Handle)
	r.POST("/avatar", avatar.Handle)
	if config.Config.SwaggerConfig.SwaggerEnabled {
		r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.Run(fmt.Sprintf(":%d", config.Config.ListenPort))
}
