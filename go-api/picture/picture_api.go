package picture

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/common"
	"github.com/peashoot/peapi/config"
	"github.com/peashoot/peapi/filemanager"
	uuid "github.com/satori/go.uuid"
)

type fillWorddRequestDto struct {
	OriImgURL string `json:"oriImgUrl"` // 原始图片url
	FillWords string `json:"fillWords"` // 填充字符串
	Font      string `json:"font"`      // 字体类型
	Scale     int    `json:"scale"`     // 图片缩放比例
	Step      int    `json:"step"`      // 字符间距
}

type fillWordsResponseDto struct {
	common.BaseResponse
	ResultImgURL string `json:"resultImgUrl"` // 生成后的图片url
}

// FillImgWithWords 用字符串填充图片
func FillImgWithWords(c *gin.Context) {
	resp := fillWordsResponseDto{}
	resp.Code = http.StatusBadRequest
	resp.Message = "Bad request"
	defer func() {
		if err := recover(); err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = "Internal server error"
			log.Println("unexpected exception appeared, the detail of error is", err)
			c.JSON(http.StatusInternalServerError, resp)
		}
	}()
	req := fillWorddRequestDto{}
	if err := c.BindJSON(&req); err != nil {
		panic(err)
	}
	checkImgURL(req.OriImgURL)
	suffix := req.OriImgURL[strings.LastIndexByte(req.OriImgURL, '.')+1:]
	if common.IndexOf([]string{"png", "jpg", "gif", "bmp"}, suffix) < 0 {
		panic("invalid original image suffix")
	}
	var oriImgPath string
	prefix := strings.TrimSuffix(config.Config.HostName, "/") + "/" +
		strings.Trim(config.Config.FileUploadConfig.FileNetURLPrefix, "/")
	if strings.HasPrefix(req.OriImgURL, prefix) {
		// 如果图片前缀是本地文件上传网络路径前缀，直接获取本地文件地址
		oriImgPath = strings.ReplaceAll(req.OriImgURL, prefix, config.Config.FileUploadConfig.FileStoreFolder)
	} else if strings.HasPrefix(req.OriImgURL, "http://") || strings.HasPrefix(req.OriImgURL, "https://") {
		var err error
		// 如果是网络路径，先下载到本地
		if _, oriImgPath, err = filemanager.DownloadFile(req.OriImgURL, "image"); err != nil {
			panic(err)
		}
	} else {
		panic("invalid original image url")
	}
	newImgName := strings.ReplaceAll(uuid.NewV4().String(), "-", "") + "." + suffix
	newImgPath := strings.TrimSuffix(config.Config.PictureOperateConfig.PictureGenerateFolder, string(os.PathSeparator)) +
		string(os.PathSeparator) + newImgName
	if err := FillWordsIntoPic(oriImgPath, newImgPath, req.FillWords, config.Config.PictureOperateConfig.PictureGenerateFonts[req.Font], req.Scale, req.Step); err != nil {
		panic(err)
	}
	newImgURL := strings.TrimSuffix(config.Config.HostName, "/") + "/" +
		strings.Trim(config.Config.PictureOperateConfig.PictureGenerateURLPrefix, "/") + "/" +
		newImgName
	resp.Code = http.StatusOK
	resp.Message = "OK"
	resp.ResultImgURL = newImgURL
	c.JSON(http.StatusOK, resp)
}

// checkImgURL 检查图片url
func checkImgURL(imgURL string) {
	if strings.Contains(imgURL, "..") || strings.Contains(imgURL, "./") || strings.Contains(imgURL, "/.") {
		panic("unvalid image url")
	}
}
