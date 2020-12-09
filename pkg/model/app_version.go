package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"strings"
)

type AppVersion struct {
	Id                   int    `json:"id" gorm:"column:id"`
	AppName              string `json:"app_name" gorm:"column:app_name"`                               //名称
	AppUrl               string `json:"app_url" gorm:"column:app_url"`                                 //地址
	AppVersionName       string `json:"app_version_name" gorm:"column:app_version_name"`               //版本名称
	AppVersionCode       string `json:"app_version_code" gorm:"column:app_version_code"`               //版本 Code
	IsForceUserUpdate    int    `json:"is_force_user_update" gorm:"column:is_force_user_update"`       //是否强制用户更新    0:否 1:是
	IsOptionalUserUpdate int    `json:"is_optional_user_update" gorm:"column:is_optional_user_update"` //是否用户选择更新  0:否 1:是
	IsSilentUpdate       int    `json:"is_silent_update" gorm:"column:is_silent_update"`               //是否静默更新  0:否 1:是
	IsHotUpdate          int    `json:"is_hot_update" gorm:"column:is_hot_update"`                     //是否热更新  0:否 1:是
	CreateTime           string `json:"create_time" gorm:"column:create_time"`                         //创建时间
	CreateUserId         int    `json:"create_user_id" gorm:"column:create_user_id"`                   //创建人
	UpdateTime           string `json:"update_time" gorm:"column:update_time"`                         //最后修改时间
	UpdateUserId         int    `json:"update_user_id" gorm:"column:update_user_id"`                   //最后修改人
}

type AppParam struct {
	DeviceNo           string `json:"device_no" `
	DeviceModel        string `json:"device_model" `
	ChannelCode        string `json:"channel_code" `
	Resolution         string `json:"resolution" `
	Language           string `json:"language" `
	StorageSpace       string `json:"storage_space" `
	AndroidVersionName string `json:"android_version_name" `
	AndroidVersionCode string `json:"android_version_code" `
	SdkName            string `json:"sdk_name" `
	SdkVersionName     string `json:"sdk_version_name" `
	SdkVersionCode     string `json:"sdk_version_code" `
	AppName            string `json:"app_name" `
	AppVersionName     string `json:"app_version_name" `
	AppVersionCode     string `json:"app_version_code" `
	CpuArch            string `json:"cpu_arch" `
}

func (AppVersion) TableName() string {
	return "app_version"
}

// 获取最新 APP 版本
func GetNewAppVersion(data map[string]interface{}) (*AppVersion, error) {
	appParam := AppParam{
		DeviceNo:           data["device_no"].(string),
		DeviceModel:        data["device_model"].(string),
		ChannelCode:        data["channel_code"].(string),
		Resolution:         data["resolution"].(string),
		Language:           data["language"].(string),
		StorageSpace:       data["storage_space"].(string),
		AndroidVersionName: data["android_version_name"].(string),
		AndroidVersionCode: data["android_version_code"].(string),
		SdkName:            data["sdk_name"].(string),
		SdkVersionName:     data["sdk_version_name"].(string),
		SdkVersionCode:     data["sdk_version_code"].(string),
		AppName:            data["app_name"].(string),
		AppVersionName:     data["app_version_name"].(string),
		AppVersionCode:     data["app_version_code"].(string),
		CpuArch:            data["cpu_arch"].(string),
	}

	appVersion, err := FetchApp(appParam)
	if err != nil {
		return nil, err
	}

	return appVersion, nil
}

func FetchApp(ap AppParam) (*AppVersion, error) {
	var appVersion *AppVersion

	//取所有有效版本
	var allVersion []*AppVersion
	err := trailerDb.Order("update_time DESC").Find(&allVersion).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, version := range allVersion {
		blackList, err := GetBlackByAppId(version.Id) //获取该版本对应的黑白名单规则
		if err != nil {
			return nil, err
		}
		if CheckBlackRole(ap, blackList) { //验证黑名单规则
			return version, nil
		}
	}

	return appVersion, nil
}

// 比对参数与黑名单规则
func CheckBlackRole(ap AppParam, blackList []*AppBlackWhite) bool {
	for _, v := range blackList {
		roleJson := v.BlackList
		var ab *AppBlackList
		json.Unmarshal([]byte(roleJson), &ab)

		if len(ap.DeviceNo) > 0 {
			if strings.Contains(ab.DeviceNo, ap.DeviceNo) { //  设备号
				return false
			}
		}
		if len(ap.DeviceModel) > 0 {
			if strings.Contains(ab.DeviceModel, ap.DeviceModel) { //设备型号
				return false
			}
		}
		if len(ap.ChannelCode) > 0 {
			if strings.Contains(ab.ChannelCode, ap.ChannelCode) { //渠道码
				return false
			}
		}
		if len(ap.Resolution) > 0 {
			if strings.Contains(ab.Resolution, ap.Resolution) { //分辨率
				return false
			}
		}
		if len(ap.Language) > 0 {
			if strings.Contains(ab.Language, ap.Language) { //语言
				return false
			}
		}
		if len(ap.StorageSpace) > 0 {
			if strings.Contains(ab.StorageSpace, ap.StorageSpace) { //存储空间
				return false
			}
		}
		if len(ap.AndroidVersionName) > 0 {
			if strings.Contains(ab.AndroidVersionName, ap.AndroidVersionName) { //Android 名称
				return false
			}
		}
		if len(ap.AndroidVersionCode) > 0 {
			if strings.Contains(ab.AndroidVersionCode, ap.AndroidVersionCode) { //Android Code
				return false
			}
		}
		if len(ap.SdkName) > 0 {
			if strings.Contains(ab.SdkName, ap.SdkName) { //Sdk 名称
				return false
			}
		}
		if len(ap.SdkVersionName) > 0 {
			if strings.Contains(ab.SdkVersionName, ap.SdkVersionName) { //Sdk 版本名称
				return false
			}
		}
		if len(ap.SdkVersionCode) > 0 {
			if strings.Contains(ab.SdkVersionCode, ap.SdkVersionCode) { //Sdk 版本 Code
				return false
			}
		}
		if len(ap.AppName) > 0 {
			if strings.Contains(ab.AppName, ap.AppName) { //App name
				return false
			}
		}
		if len(ap.AppVersionName) > 0 {
			if strings.Contains(ab.AppVersionName, ap.AppVersionName) { //App 版本名称
				return false
			}
		}
		if len(ap.AppVersionCode) > 0 {
			if strings.Contains(ab.AppVersionCode, ap.AppVersionCode) { //App 版本 Code
				return false
			}
		}
		if len(ap.CpuArch) > 0 {
			if strings.Contains(ab.CpuArch, ap.CpuArch) { //CPU 架构
				return false
			}
		}
	}

	return true
}
