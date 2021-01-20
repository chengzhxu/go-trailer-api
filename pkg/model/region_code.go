package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"os"
)

type RegionInfo struct {
	RegionCode string `json:"region_code" gorm:"column:region_code"` //region_code
	RegionName string `json:"region_name" gorm:"column:region_name"` //region_name
}

var filePath = "conf/t_region_code.txt"

type RegionMap map[string]string

func (RegionInfo) TableName() string {
	return "t_region_code"
}

//获取所有的 Region
func GetAllRegionInfo() ([]*RegionInfo, error) {
	var allRegion []*RegionInfo
	err := trailerDb.Find(&allRegion).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return allRegion, nil
}

//将 Region 信息写入文件
func WriteRegionToFile() {
	_, err := os.Stat(filePath)
	if err != nil {
		os.Create(filePath)
	} else {
		f, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
		defer f.Close()
	}

	allRegion, _ := GetAllRegionInfo()

	regionMap := make(RegionMap)
	if len(allRegion) > 0 {
		for _, r := range allRegion {
			//fmt.Printf("%v", r.RegionName)
			//os.Exit(0)
			if r.RegionCode != "" && r.RegionName != "" {
				regionMap[r.RegionName] = r.RegionCode
			}

		}
	}

	data, err := json.MarshalIndent(regionMap, "", " ")
	if err != nil {
		panic(err)
	}

	//将数据data写入文件filePath中，并且修改文件权限为755
	ioutil.WriteFile(filePath, data, 0755)
}

func ReadRegionFromFile(regionName string) string {
	if regionName == "" {
		return ""
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}

	regionMap := make(RegionMap)
	//Unmarshal将data数据转换成指定的结构体类型,经过json转换后data中的数据已写入map中
	if err := json.Unmarshal(data, regionMap); err != nil {
		return ""
	}

	return regionMap[regionName]
}
