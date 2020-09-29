package cnarea

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/database"
)

// queryConditionDTO 请求数据传输实体
type queryConditionDTO struct {
	ParentCode   int64  `json:"parent_code"`   // 父级区域代码
	NameContains string `json:"name_contains"` // 名称字串
	PageIndex    int    `json:"page_index"`    // 分页页数 从1开始
	PageSize     int    `json:"page_size"`     // 分页大小
	ZipCode      string `json:"zip_code"`      // 邮政编码
}

// queryResultDTO 返回数据传输实体
type queryResultDTO struct {
	Code       int    `json:"code"`        // 返回代码 200 成功
	Message    string `json:"message"`     // 错误说明
	TotalCount int    `json:"total_count"` // 当前记录数
	PageIndex  int    `json:"page_index"`  // 当前页码 从1开始
	PageSize   int    `json:"page_size"`   // 分页大小
	Data       []dto  `json:"data"`        // 数据
}

// DTO 区域数据传输实体
type dto struct {
	AreaCode   int64   `json:"area_code"`   // 行政代码
	ZipCode    int     `json:"zip_code"`    // 邮政编码
	CityCode   string  `json:"city_code"`   // 区号
	Name       string  `json:"name"`        // 名称
	ShortName  string  `json:"short_name"`  // 简称
	MergerName string  `json:"merger_name"` // 组合名
	PinYin     string  `json:"pinyin"`      // 拼音
	Longitude  float64 `json:"longitude"`   // 经度
	Latitude   float64 `json:"latitude"`    // 纬度
}

// Handle 处理请求
// @Summary 查询区域信息
// @Description 通过父级区域代码、邮政编码等查询条件分页查询符合条件的区域信息
// @Param condition body queryConditionDTO true "查询条件"
// @Accept json
// @Produce json
// @Success 200 {object} queryResultDTO
// @Router /cnarea [post]
func Handle(c *gin.Context) {
	resp := &queryResultDTO{
		Code:    400,
		Message: "Bad Request",
	}
	req := &queryConditionDTO{}
	if err := c.BindJSON(req); err != nil {
		log.Println("Failure to read json body, the detail of error is", err)
		c.JSON(http.StatusBadRequest, resp)
	}
	resp.Code = 500
	resp.Message = "internal server error"
	resp.PageIndex = req.PageIndex
	resp.PageSize = req.PageSize
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
	}
	var areas []DO
	defer db.Close()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("unexpected exception appeared, the detail of error is", err)
			c.JSON(http.StatusInternalServerError, resp)
		}
	}()
	db = db.Where("parent_code = ?", req.ParentCode)
	if req.NameContains != "" {
		db = db.Where("name like ?", "%"+req.NameContains+"%")
	}
	if req.ZipCode != "" {
		db = db.Where("zip_code = ?", req.ZipCode)
	}
	var totalCount int
	db.Model(&DO{}).Count(&totalCount)
	db.Offset(req.PageSize * (req.PageIndex - 1)).Limit(req.PageSize).Find(&areas)
	resp.Code = 200
	resp.Message = "success"
	resp.TotalCount = totalCount
	for _, area := range areas {
		entity := &dto{
			AreaCode:   area.AreaCode,
			ZipCode:    area.ZipCode,
			CityCode:   area.CityCode,
			Name:       area.Name,
			ShortName:  area.ShortName,
			MergerName: area.MergerName,
			PinYin:     area.PinYin,
			Longitude:  area.Longitude,
			Latitude:   area.Latitude,
		}
		resp.Data = append(resp.Data, *entity)
	}
	c.JSON(http.StatusOK, resp)
}
