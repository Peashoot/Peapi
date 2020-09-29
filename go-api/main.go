package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/peashoot/peapi/avatar"
	"github.com/peashoot/peapi/common"
	"github.com/peashoot/peapi/config"
	_ "github.com/peashoot/peapi/docs"
	"github.com/peashoot/peapi/router"
)

func init() {
	flag.StringVar(&confPath, "c", "conf.json", "path of configuration file")
	flag.Parse()
}

var (
	confPath string // 配置文件路径
)

func main() {
	config.ReadConfig(confPath)
	if !config.Config.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	router.RegisterRouter(r)
	startTimer()
	r.Run(fmt.Sprintf(":%d", config.Config.ListenPort))
}

func startTimer() *common.Crontab {
	crontab := common.NewCrontab()
	if err := crontab.AddByFunc("clear avatar", config.Config.AvatarConfig.AvatarCleanEventCron, avatar.ClearOverdueFiles); err != nil {
		log.Print("fail to start timer of clear avatar, err is:", err)
	}
	crontab.Start()
	return crontab
}
