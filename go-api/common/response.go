package common

// BaseResponse 基本响应
type BaseResponse struct {
	Code    int    `json:"code"`    // 响应代码
	Message string `json:"message"` // 错误说明
}
