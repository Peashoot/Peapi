package filemanager

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/common"
	"github.com/peashoot/peapi/config"
	uuid "github.com/satori/go.uuid"
)

// uploadResponseDto 上传文件响应
type uploadResponseDto struct {
	common.BaseResponse
	RemoteFileURL  string `json:"remoteFileUrl"`  // 上传文件链接
	RenameFileName string `json:"renameFileName"` // 重命名后文件名称
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	resp := uploadResponseDto{}
	resp.Code = http.StatusBadRequest
	resp.Message = "Bad request"
	defer func() {
		if err := recover(); err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = "Internal server error"
			log.Println("an unexpected exception appeared, the detail of error is", err)
			c.JSON(http.StatusInternalServerError, resp)
		}
	}()
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}
	filetype := form.Value["filetype"][0]
	if filetype == "" {
		filetype = "text"
	}
	var filename, oriFilename string
	if filenameVal := form.Value["filename"]; filenameVal != nil && len(filenameVal) > 0 {
		oriFilename = filenameVal[0]
		filename = oriFilename
	}
	if filename == "" {
		filename = strings.Replace(uuid.NewV4().String(), "-", "", -1) + "." + getDefaultFileTypeSuffix(filetype)
	} else if suffix := getFileSuffix(oriFilename); checkFileSuffix(filetype, suffix) {
		filename = strings.Replace(uuid.NewV4().String(), "-", "", -1) + "." + suffix
	} else {
		panic("data verification failed")
	}
	localFolderPath := strings.TrimSuffix(config.Config.FileUploadConfig.FileStoreFolder, string(os.PathSeparator)) + string(os.PathSeparator) + filetype
	os.MkdirAll(localFolderPath, os.ModePerm)
	localFilePath := localFolderPath + string(os.PathSeparator) + filename
	fileHeader := form.File["file"][0]
	file, err := fileHeader.Open()
	if err != nil {
		panic(err)
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	localFile, err := os.Create(localFilePath)
	if err != nil {
		panic(err)
	}
	if _, err = localFile.Write(fileBytes); err != nil {
		panic(err)
	}
	resp.Code = http.StatusOK
	resp.Message = "OK"
	resp.RemoteFileURL = strings.TrimSuffix(config.Config.HostName, "/") + "/" +
		strings.Trim(config.Config.FileUploadConfig.FileNetURLPrefix, "/") + "/" +
		filetype + "/" +
		filename
	resp.RenameFileName = filename
	c.JSON(http.StatusOK, resp)
}

// DownloadFile 下载网络图片到本地
func DownloadFile(fileURL string, filetype string) (renameFilename string, localFilePath string, retErr error) {
	if suffix := getFileSuffix(fileURL); checkFileSuffix(filetype, suffix) {
		renameFilename = strings.ReplaceAll(uuid.NewV4().String(), "-", "") + "." + suffix
	} else {
		retErr = errors.New("unmatched file suffix")
		return
	}
	req, err := http.NewRequest(http.MethodGet, fileURL, nil)
	if err != nil {
		retErr = err
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		retErr = err
		return
	}
	defer res.Body.Close()
	fileBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		retErr = err
		return
	}
	localFolderPath := strings.TrimSuffix(config.Config.FileUploadConfig.FileStoreFolder, string(os.PathSeparator)) + string(os.PathSeparator) + filetype
	os.MkdirAll(localFolderPath, os.ModePerm)
	localFilePath = localFolderPath + string(os.PathSeparator) + renameFilename
	localFile, err := os.Create(localFilePath)
	if err != nil {
		retErr = err
		return
	}
	_, err = localFile.Write(fileBytes)
	retErr = err
	return
}

// getDefaultFileTypeSuffix 获取具体文件类型的默认后缀名
func getDefaultFileTypeSuffix(filetype string) string {
	switch filetype {
	case "text":
		return "txt"
	case "image":
		return "jpg"
	case "markdown":
		return "md"
	}
	panic("unknown filetype")
}

// getFileSuffix 获取文件后缀
func getFileSuffix(filename string) string {
	index := strings.LastIndexByte(filename, '.')
	if index < 0 {
		panic("fail to capture suffix")
	}
	return filename[index+1:]
}

// checkFileSuffix 检查文件后缀
func checkFileSuffix(filetype, suffix string) bool {
	switch filetype {
	case "text":
		return common.IndexOf([]string{"txt", "tmp", "log", "tmp"}, suffix) >= 0
	case "image":
		return common.IndexOf([]string{"jpg", "png", "bmp", "gif", "tiff"}, suffix) >= 0
	case "markdown":
		return "md" == suffix
	}
	panic("unknown filetype")
}
