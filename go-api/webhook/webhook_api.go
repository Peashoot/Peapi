package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/config"
)

type hookResultDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Handle 处理Webhook
func Handle(c *gin.Context) {
	resp := &hookResultDTO{
		Code:    400,
		Message: "Bad request",
	}
	//我们这里只接收json的请求体
	if contentType := c.GetHeader("content-type"); contentType != "application/json" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// 首先应该对secret进行校验
	signature := c.GetHeader("X-Hub-Signature")
	if len(signature) <= 0 {
		resp.Code = 401
		resp.Message = "Unauthorized"
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	// 读取body
	bytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// 验证签名
	hash := hmac.New(sha1.New, []byte(config.Config.WebhookConfig.WebHookSecret))
	hash.Write(bytes)
	sign := "sha1=" + hex.EncodeToString(hash.Sum(nil))
	if strings.Compare(sign, signature) != 0 {
		resp.Code = 401
		resp.Message = "Unauthorized"
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	// 绑定body到json
	var body map[string]interface{}
	if err := c.BindJSON(body); err != nil {
		resp.Code = 406
		resp.Message = "Not accrptable"
		c.JSON(http.StatusNotAcceptable, resp)
		return
	}
	//body里就已经获取到我们想要的数据了, 在这里我们只拿取push的分支，更多参数参考可以参考文档
	ref, ok := body["ref"].(string)
	if !ok {
		resp.Code = 500
		resp.Message = "Internal server error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//拿到的ref：refs/heads/master，我们只需要把前面截掉，就是我们想要的分支名
	branch := strings.TrimLeft(ref, "refs/heads/")
	//我们这里只对master分支进行自动部署
	if branch != "master" {
		resp.Code = 202
		resp.Message = "Accepted"
		c.JSON(http.StatusAccepted, resp)
		return
	}
	//我这里的部署是运行一个shell脚本，我们可以在这里写自动更新和部署的逻辑
	if err := exec.Command(config.Config.WebhookConfig.WebHookShellPath).Run(); err != nil {
		resp.Code = 500
		resp.Message = "Internal server error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Code = 200
	resp.Message = "OK"
	c.JSON(http.StatusOK, resp)
}
