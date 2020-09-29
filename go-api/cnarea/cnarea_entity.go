package cnarea

// DO 中国区域实体
type DO struct {
	RecordID   int     `gorm:"column:id; PRIMARY_KEY"` // 记录id
	Level      int     `gorm:"column:level"`           // 层级
	ParentCode int64   `gorm:"column:parent_code"`     // 父级行政代码
	AreaCode   int64   `gorm:"column:area_code"`       // 行政代码
	ZipCode    int     `gorm:"column:zip_code"`        // 邮政编码
	CityCode   string  `gorm:"column:city_code"`       // 区号
	Name       string  `gorm:"column:name"`            // 名称
	ShortName  string  `gorm:"column:short_name"`      // 简称
	MergerName string  `gorm:"column:merger_name"`     // 组合名
	PinYin     string  `gorm:"column:pinyin"`          // 拼音
	Longitude  float64 `gorm:"column:lng"`             // 经度
	Latitude   float64 `gorm:"column:lat"`             // 纬度
}

// TableName 表名
func (area DO) TableName() string {
	return "cnarea_2019"
}
