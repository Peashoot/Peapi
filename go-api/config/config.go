package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// RunConfig 运行参数
type RunConfig struct {
	HostName             string               `json:"hostname"`    // 服务主机名
	DebugMode            bool                 `json:"debug_mode"`  // 调试模式
	ListenPort           int                  `json:"listen_port"` // 监听端口
	DatabaseConfig       DatabaseConfig       `json:"database"`    // 数据库配置参数
	WebhookConfig        WebhookConfig        `json:"webhook"`     // webhook配置参数
	SwaggerConfig        SwaggerConfig        `json:"swagger"`     // swagger配置参数
	AvatarConfig         AvatarConfig         `json:"avatar"`      // 头像api配置参数
	DownloadConfig       DownloadConfig       `json:"download"`    // 离线下载配置参数
	FileUploadConfig     FileUploadConfig     `json:"file_upload"` // 文件上传配置参数
	PictureOperateConfig PictureOperateConfig `json:"pic_operate"` // 图片操作配置参数
}

// DatabaseConfig 数据库配置参数
type DatabaseConfig struct {
	DBUsername string `json:"username"` // 数据库连接用户名
	DBPassword string `json:"password"` // 数据库连接密码
	DBHost     string `json:"host"`     // 数据库地址
	DBName     string `json:"name"`     // 数据库名称
}

// WebhookConfig webhook配置参数
type WebhookConfig struct {
	WebHookSecret    string `json:"secret"` // Webhook密钥
	WebHookShellPath string `json:"shell"`  // Webhook执行脚本路径
}

// SwaggerConfig swagger配置参数
type SwaggerConfig struct {
	SwaggerEnabled bool `json:"enabled"` // 是否启用Swagger文
}

// AvatarConfig 头像api配置参数
type AvatarConfig struct {
	AvatarFileFolderPath string `json:"local_path"`    // 头像文件夹本地路径
	AvatarFileNetURL     string `json:"net_url"`       // 头像网络路径
	AvatarSaveDuration   int    `json:"save_duration"` // 头像文件缓存时间（单位：分钟）
	AvatarCleanEventCron string `json:"clean_cron"`    // 头像定时清理任务轮询cron规则
}

// DownloadConfig 离线下载配置参数
type DownloadConfig struct {
	Aria2RpcURL string `json:"aria2_rpc"`    // aria2 RPC 请求地址
	Aria2Secret string `json:"aria2_secret"` // aria2 RPC 密钥
}

// FileUploadConfig 文件上传配置参数
type FileUploadConfig struct {
	FileStoreFolder  string `json:"file_store_folder"` // 文件存储文件夹路径
	FileNetURLPrefix string `json:"file_net_url"`      // 文件网络地址前缀
}

// PictureOperateConfig 图片操作配置参数
type PictureOperateConfig struct {
	PictureGenerateFolder    string            `json:"pic_generate_folder"`   // 生成的图片的储存路径
	PictureGenerateURLPrefix string            `json:"pic_generate_net_url"`  // 生成的图片的网络地址前缀
	PictureGenerateFonts     map[string]string `json:"pic_generete_font_dic"` // 图片生成过程中需要的字典
	PictureSaveDuration      int               `json:"save_duration"`         // 生成的文件缓存时间（单位：分钟）
	PictureCleanEventCron    string            `json:"clean_cron"`            // 生成图片定时清理任务轮询cron规则
}

// Config 配置项
var Config *RunConfig

// DefaultConfig 获取默认配置
func DefaultConfig() *RunConfig {
	return &RunConfig{
		HostName:   "http://127.0.0.1:8080",
		DebugMode:  false,
		ListenPort: 8080,
		DatabaseConfig: DatabaseConfig{
			DBHost: "127.0.0.1",
			DBName: "peapi",
		},
		WebhookConfig: WebhookConfig{
			WebHookSecret:    "PeashootWithGo",
			WebHookShellPath: "/usr/share/webhook.sh",
		},
		SwaggerConfig: SwaggerConfig{},
		AvatarConfig: AvatarConfig{
			AvatarCleanEventCron: "0/10 * * * * ?",
			AvatarFileFolderPath: "/usr/share/avatars",
			AvatarFileNetURL:     "/files/avatar",
			AvatarSaveDuration:   15,
		},
		DownloadConfig: DownloadConfig{
			Aria2RpcURL: "http://127.0.0.1:6800/jsonrpc",
			Aria2Secret: "",
		},
		FileUploadConfig: FileUploadConfig{
			FileStoreFolder:  "/usr/share/files",
			FileNetURLPrefix: "/files/upload",
		},
		PictureOperateConfig: PictureOperateConfig{
			PictureGenerateFolder:    "/usr/share/pictures",
			PictureGenerateURLPrefix: "/files/picture",
			PictureGenerateFonts:     make(map[string]string),
			PictureCleanEventCron:    "0/10 * * * * ?",
			PictureSaveDuration:      15,
		},
	}
}

// ReadConfig 读取配置文件
func ReadConfig(jsonPath string) {
	Config = DefaultConfig()
	jsonBytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Println("Failure to read configuration file. Detail of error is", err)
		return
	}
	if err = json.Unmarshal(jsonBytes, Config); err != nil {
		log.Println("Failure to unmarshal json configuration file. Detail of error is", err)
		return
	}
	if err := os.MkdirAll(Config.AvatarConfig.AvatarFileFolderPath, os.ModePerm); err != nil {
		log.Println("Failure to create avatar file folder. Detail of error is", err)
	}
	if err := os.MkdirAll(Config.FileUploadConfig.FileStoreFolder, os.ModePerm); err != nil {
		log.Println("Failure to create upload file folder. Detail of error is", err)
	}
	if err := os.MkdirAll(Config.PictureOperateConfig.PictureGenerateFolder, os.ModePerm); err != nil {
		log.Println("Failure to create generate picture folder. Detail of error is", err)
	}
	dir, _ := os.Getwd()
	if Config.PictureOperateConfig.PictureGenerateFonts["consola"] == "" {
		Config.PictureOperateConfig.PictureGenerateFonts["consola"] = path.Join(dir, "scripts", "source", "consola.ttf")
	}
	if Config.PictureOperateConfig.PictureGenerateFonts["simsun"] == "" {
		Config.PictureOperateConfig.PictureGenerateFonts["simsun"] = path.Join(dir, "scripts", "source", "simsun.ttf")
	}
	return
}
