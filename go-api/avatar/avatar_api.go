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

type requestDTO struct {
	Sex      string `json:"sex"`      // 性别
	Username string `json:"username"` // 用户名
}

type responseDTO struct {
	Code    int    `json:"code"`    // 响应代码
	Message string `json:"message"` // 错误说明
	ImgURL  string `json:"imgUrl"`  // 图片URL
}

// Handle 处理头像生成
func Handle(c *gin.Context) {
	response := &responseDTO{
		Code:    500,
		Message: "Internal Server Error",
	}
	request := &requestDTO{}
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
