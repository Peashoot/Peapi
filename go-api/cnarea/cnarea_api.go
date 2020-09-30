package cnarea

import (
	"errors"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/peashoot/peapi/common"
	"github.com/peashoot/peapi/database"
)

// listRequestDTO 请求数据传输实体
type listRequestDTO struct {
	ParentCode   int64  `json:"parent_code"`   // 父级区域代码
	NameContains string `json:"name_contains"` // 名称字串
	PageIndex    int    `json:"page_index"`    // 分页页数 从1开始
	PageSize     int    `json:"page_size"`     // 分页大小
	ZipCode      string `json:"zip_code"`      // 邮政编码
}

// listResponseDTO 返回数据传输实体
type listResponseDTO struct {
	common.BaseResponse
	TotalCount int       `json:"total_count"` // 当前记录数
	PageIndex  int       `json:"page_index"`  // 当前页码 从1开始
	PageSize   int       `json:"page_size"`   // 分页大小
	Data       []areaDTO `json:"data"`        // 数据
}

// DTO 区域数据传输实体
type areaDTO struct {
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

// areaWithParentDTO 传输实体附带父级实体
type areaWithParentDTO struct {
	areaDTO
	Parent *areaWithParentDTO `json:"parent"` // 父节点
}

// transformDOToDTO 转换数据库实体为传输实体
func transformDOToDTO(area DO) areaDTO {
	return areaDTO{
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
}

// transformDOToDTOWithParant 转换数据库实体为传输实体并带上父级实体
func transformDOToDTOWithParant(area DO) (areaWithParentDTO, error) {
	inner := transformDOToDTO(area)
	if area.ParentCode == 0 {
		return areaWithParentDTO{
			areaDTO: inner,
		}, nil
	}
	var err error
	area, err = getAreaByCode(area.ParentCode)
	if err != nil {
		return areaWithParentDTO{}, err
	}
	var parent areaWithParentDTO
	parent, err = transformDOToDTOWithParant(area)
	return areaWithParentDTO{
		areaDTO: inner,
		Parent:  &parent,
	}, nil
}

// List 条件查询区域信息
// @Summary 查询区域信息
// @Description 通过父级区域代码、邮政编码等查询条件分页查询符合条件的区域信息
// @Param condition body queryConditionDTO true "查询条件"
// @Accept json
// @Produce json
// @Success 200 {object} queryResultDTO
// @Router /cnarea [post]
func List(c *gin.Context) {
	resp := &listResponseDTO{}
	resp.Code = 400
	resp.Message = "Bad Request"
	req := &listRequestDTO{}
	if err := c.BindJSON(req); err != nil {
		log.Println("Failure to read json body, the detail of error is", err)
		c.JSON(http.StatusBadRequest, resp)
	}
	resp.Code = 500
	resp.Message = "Internal server error"
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
	resp.Message = "OK"
	resp.TotalCount = totalCount
	for _, area := range areas {
		entity := transformDOToDTO(area)
		resp.Data = append(resp.Data, entity)
	}
	c.JSON(http.StatusOK, resp)
}

// locateRequestDTO 定位请求实体
type locateRequestDTO struct {
	Longtitude float64 `json:"longtitude"` // 经度
	Latitude   float64 `json:"latitude"`   // 纬度
}

// locateResponseDTO 定位响应实体
type locateResponseDTO struct {
	common.BaseResponse
	Area areaWithParentDTO `json:"area"` // 区域信息
}

// Locate 定位
func Locate(c *gin.Context) {
	resp := locateResponseDTO{}
	resp.Code = 400
	resp.Message = "Bad Request"
	req := &locateRequestDTO{}
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
	area, err := queryAreaEntity(req.Longtitude, req.Latitude)
	if err != nil {
		panic(err)
	}
	resp.Code = 200
	resp.Message = "OK"
	resp.Area, err = transformDOToDTOWithParant(area)
	c.JSON(http.StatusOK, resp)
}

// queryAreaEntity 根据经纬度获取区域编号
func queryAreaEntity(longitude float64, latitude float64) (DO, error) {
	var areas []DO
	var err error
	accuracyArray := []float64{0.01, 0.02, 0.03, 0.05, 0.1, 0.2, 0.3, 0.5, 1, 2, 3, 5, 10}
	for _, accuracy := range accuracyArray {
		areas, err = getNearByArea(longitude, latitude, accuracy)
		if err != nil {
			return DO{}, err
		}
		if len(areas) > 0 {
			break
		}
	}
	if len(areas) == 0 {
		return DO{}, errors.New("can't find matched area")
	}
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].GetDistance(longitude, latitude) < areas[j].GetDistance(longitude, latitude)
	})
	return areas[0], nil
}

// getNearByArea 根据精度搜索范围
func getNearByArea(longitude float64, latitude float64, accuracy float64) ([]DO, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var areas []DO
	db = db.Where("level = ?", 4).
		Where("lng > ? and lng < ?", longitude-accuracy, longitude+accuracy).
		Where("lat > ? and lat < ?", latitude-accuracy, latitude+accuracy).
		Find(&areas)
	return areas, nil
}

// getAreaByCode 根据区域编号获取区域信息
func getAreaByCode(areaCode int64) (DO, error) {
	db, err := database.Connect()
	if err != nil {
		return DO{}, err
	}
	defer db.Close()
	var areas []DO
	db = db.Where("area_code = ?", areaCode).
		Limit(1).
		Find(&areas)
	if len(areas) < 1 {
		return DO{}, errors.New("no matched record")
	}
	return areas[0], nil
}
