package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// RunConfig 运行参数
type RunConfig struct {
	DebugMode      bool           `json:"debug_mode"`  // 调试模式
	ListenPort     int            `json:"listen_port"` // 监听端口
	DatabaseConfig DatabaseConfig `json:"database"`    // 数据库配置参数
	WebhookConfig  WebhookConfig  `json:"webhook"`     // webhook配置参数
	SwaggerConfig  SwaggerConfig  `json:"swagger"`     // swagger配置参数
	AvatarConfig   AvatarConfig   `json:"avatar"`      // 头像api配置参数
	DownloadConfig DownloadConfig `json:"download"`    // 离线下载配置参数
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

// DownloadConfig 离线下载配置文件
type DownloadConfig struct {
	Aria2RpcURL string `json:"aria2_rpc"`    // aria2 RPC 请求地址
	Aria2Secret string `json:"aria2_secret"` // aria2 RPC 密钥
}

// Config 配置项
var Config *RunConfig

// DefaultConfig 获取默认配置
func DefaultConfig() *RunConfig {
	return &RunConfig{
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
			AvatarFileNetURL:     "http://127.0.0.1:8089",
			AvatarSaveDuration:   15,
		},
		DownloadConfig: DownloadConfig{
			Aria2RpcURL: "http://127.0.0.1:6800/jsonrpc",
			Aria2Secret: "",
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
	return
}
