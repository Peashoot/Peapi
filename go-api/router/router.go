package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/avatar"
	"github.com/peashoot/peapi/cnarea"
	"github.com/peashoot/peapi/config"
	"github.com/peashoot/peapi/filemanager"
	"github.com/peashoot/peapi/picture"
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
	engine.POST("/file/upload", filemanager.UploadFile)
	engine.POST("/picture/fillwords", picture.FillImgWithWords)
	engine.Static(config.Config.AvatarConfig.AvatarFileNetURL, config.Config.AvatarConfig.AvatarFileFolderPath)
	engine.Static(config.Config.FileUploadConfig.FileNetURLPrefix, config.Config.FileUploadConfig.FileStoreFolder)
	engine.Static(config.Config.PictureOperateConfig.PictureGenerateURLPrefix, config.Config.PictureOperateConfig.PictureGenerateFolder)
	if config.Config.SwaggerConfig.SwaggerEnabled {
		engine.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
