package avatar

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/o1egl/govatar"
	"github.com/peashoot/peapi/common"
	"github.com/peashoot/peapi/config"
	uuid "github.com/satori/go.uuid"
)

type generateRequestDTO struct {
	Sex      string `json:"sex"`      // 性别
	Username string `json:"username"` // 用户名
}

type generateResponseDTO struct {
	common.BaseResponse
	ImgURL string `json:"imgUrl"` // 图片URL
}

// Generate 处理头像生成
func Generate(c *gin.Context) {
	response := &generateResponseDTO{}
	response.Code = 500
	response.Message = "Internal Server Error"
	request := &generateRequestDTO{}
	if err := c.BindJSON(request); err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	filePath := config.Config.AvatarConfig.AvatarFileLocalPath
	if !strings.HasSuffix(filePath, string(os.PathSeparator)) {
		filePath += string(os.PathSeparator)
	}
	fileName := strings.Replace(uuid.NewV4().String(), "-", "", -1) + ".png"
	filePath += fileName
	sex := govatar.FEMALE
	arr := []string{"man", "nan", "男"}
	if common.IndexOf(arr, strings.ToLower(request.Sex)) > 0 {
		sex = govatar.MALE
	}
	if "" == request.Username {
		if err := govatar.GenerateFile(sex, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	} else {
		if err := govatar.GenerateFileForUsername(sex, request.Username, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	}
	netURL := config.Config.AvatarConfig.AvatarFileNetURL
	if !strings.HasSuffix(netURL, "/") {
		netURL += "/"
	}
	response.Code = 200
	response.Message = "OK"
	response.ImgURL = netURL + fileName
	c.JSON(http.StatusOK, response)
}
