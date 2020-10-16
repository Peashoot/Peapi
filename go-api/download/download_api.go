package download

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/common"
)

// RegisterDownloadRouter 注册下载路由
func RegisterDownloadRouter(engine *gin.Engine) {

}

type baseDownloadRequest struct {
	Token string `json:"token"` //鉴权token
}

type createLinkMissionRequest struct {
	baseDownloadRequest
	Link string `json:"link"` // 下载链接
}

type createLinkMissionResponse struct {
	common.BaseResponse
	GID string `json:"gid"` // 任务id
}

// CreateLinkMission 创建链接任务
func CreateLinkMission(c *gin.Context) {
	resp := createLinkMissionResponse{}
	resp.Code = 400
	resp.Message = "Bad Request"
	req := createLinkMissionRequest{}
	if err := c.BindJSON(req); err != nil {
		log.Println("Failure to read json body, the detail of error is", err)
		c.JSON(http.StatusBadRequest, resp)
	}
	resp.Code = 500
	resp.Message = "Internal server error"
	defer func() {
		err := recover()
		if err != nil {
			log.Println("unexpected exception appeared, the detail of error is", err)
			c.JSON(http.StatusInternalServerError, resp)
		}
	}()
	if gid, err := AddURI(req.Link); err != nil {
		panic(err)
	} else {
		resp.Code = 200
		resp.Message = "OK"
		resp.GID = gid
		c.JSON(http.StatusOK, resp)
	}
}
