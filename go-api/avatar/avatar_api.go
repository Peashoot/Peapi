package avatar

import (
	"math/rand"
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
	ImgName string `json:"imgName"` // 图片名称
	ImgURL  string `json:"imgUrl"`  // 图片URL
}

// Generate 生成头像
// @Summary 生成头像
// @Description 输入性别用户名等信息生成卡通头像
// @Param factor body generateRequestDTO true "生成因子"
// @Accept json
// @Produce json
// @Success 200 {object} generateResponseDTO
// @Router /avatar/generate [post]
func Generate(c *gin.Context) {
	response := &generateResponseDTO{}
	response.Code = 500
	response.Message = "Internal Server Error"
	request := &generateRequestDTO{}
	if err := c.BindJSON(request); err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fileName := strings.Replace(uuid.NewV4().String(), "-", "", -1) + ".png"
	filePath := GetAvatarFilePath(fileName)
	sex := govatar.FEMALE
	if rand.Int()%2 == 1 {
		sex = govatar.MALE
	}
	arr := []string{"man", "nan", "男"}
	if common.IndexOf(arr, strings.ToLower(request.Sex)) >= 0 {
		sex = govatar.MALE
	} else if arr = []string{"woman", "female", "女"}; common.IndexOf(arr, strings.ToLower(request.Sex)) >= 0 {
		sex = govatar.FEMALE
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
	response.ImgName = fileName
	response.ImgURL = netURL + fileName
	c.JSON(http.StatusOK, response)
}

// GetAvatarFilePath 获取头像文件路径
func GetAvatarFilePath(fileName string) string {
	filePath := config.Config.AvatarConfig.AvatarFileFolderPath
	if !strings.HasSuffix(filePath, string(os.PathSeparator)) {
		filePath += string(os.PathSeparator)
	}
	filePath += fileName
	return filePath
}
