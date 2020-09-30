package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/avatar"
	"github.com/peashoot/peapi/cnarea"
	"github.com/peashoot/peapi/config"
	"github.com/peashoot/peapi/webhook"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRouter 注册路由
func RegisterRouter(engine *gin.Engine) {
	engine.POST("/cnarea/list", cnarea.List)
	engine.POST("/cnarea/locate", cnarea.Locate)
	engine.POST("/webhook/execute", webhook.Exceute)
	engine.POST("/avatar/generate", avatar.Generate)
	if config.Config.SwaggerConfig.SwaggerEnabled {
		engine.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
