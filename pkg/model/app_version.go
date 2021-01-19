package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"go-trailer-api/pkg/util"
	"strconv"
	"strings"
)

type AppVersion struct {
	Id                   int    `json:"id" gorm:"column:id"`
	AppName              string `json:"app_name" gorm:"column:app_name"`                               //名称
	AppUrl               string `json:"app_url" gorm:"column:app_url"`                                 //地址
	AppVersionName       string `json:"app_version_name" gorm:"column:app_version_name"`               //版本名称
	AppVersionCode       string `json:"app_version_code" gorm:"column:app_version_code"`               //版本 Code
	Remark               string `json:"remark" gorm:"column:remark"`                                   //版本介绍
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
	IsHotUpdate        int    `json:"is_hot_update"` //是否热更  0:否  1:是
	IsSecure           bool   `json:"isSecure"`      //判断返回链接形式 https or http;
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
		IsHotUpdate:        data["is_hot_update"].(int),
		IsSecure:           data["is_secure"].(bool),
	}

	appVersion, err := FetchApp(appParam)
	if err != nil {
		return nil, err
	}

	return appVersion, nil
}

func FetchApp(ap AppParam) (*AppVersion, error) {
	var appVersion *AppVersion

	isHotUpdate := ap.IsHotUpdate                        //是否热更版本
	appName := ap.AppName                                //版本名称
	sdkName := ap.SdkName                                //SDK 名称
	appVersionCode, _ := strconv.Atoi(ap.AppVersionCode) //版本 code
	sdkVersionCode, _ := strconv.Atoi(ap.SdkVersionCode) //版本 code
	//取所有有效版本
	var allVersion []*AppVersion

	if strings.ReplaceAll(sdkName, " ", "") != "" {
		appName = sdkName
	}
	if sdkVersionCode > 0 {
		appVersionCode = sdkVersionCode
	}
	err := trailerDb.Where("is_hot_update = ? and app_name = ? and app_version_code > ? ", isHotUpdate, appName, appVersionCode).Order("update_time DESC").Find(&allVersion).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, version := range allVersion {
		ruleList, err := GetBlackWhiteByAppId(version.Id) //获取该版本对应的黑白名单规则
		if err != nil {
			return nil, err
		}
		if CheckBlackWhiteRule(ap, ruleList) { //验证黑百名单规则
			if !ap.IsSecure {
				strings.Replace(version.AppUrl, "https://", "http://", 1)
			}
			return version, nil
		}
	}

	return appVersion, nil
}

// 比对参数与黑/白名单规则
func CheckBlackWhiteRule(ap AppParam, ruleList []*AppBlackWhite) bool {
	for _, v := range ruleList {
		blackList := v.BlackList //黑名单
		whiteList := v.WhiteList //白名单

		var black *AppBlackList //规则
		if blackList != "" {
			json.Unmarshal([]byte(blackList), &black)
		}

		var white *AppBlackList //规则
		if whiteList != "" {
			json.Unmarshal([]byte(whiteList), &white)
		}

		if blackList != "" && whiteList == "" { // 有黑名单 & 无白名单
			return checkBlack(ap, black)
		}
		if blackList == "" && whiteList != "" { // 无黑名单 & 有白名单
			return checkWhite(ap, white)
		}
		if blackList != "" && whiteList != "" { // 有黑名单 & 有白名单
			return checkBlackWhite(ap, black, white)
		}
	}

	return true
}

// 验证黑名单 - 有黑名单 & 无白名单 - 若在黑名单内，则不更新；否则更新；
func checkBlack(ap AppParam, black *AppBlackList) bool {
	if len(ap.DeviceNo) > 0 && black.DeviceNo != "" {
		dnArr := strings.Split(black.DeviceNo, ",")
		if util.StrInArray(ap.DeviceNo, dnArr) { //  设备号
			return false
		}
	}
	if len(ap.DeviceModel) > 0 && black.DeviceModel != "" {
		dmArr := strings.Split(black.DeviceModel, ",")
		if util.StrInArray(ap.DeviceModel, dmArr) { //设备型号
			return false
		}
	}
	if len(ap.ChannelCode) > 0 && black.ChannelCode != "" {
		ccArr := strings.Split(black.ChannelCode, ",")
		if util.StrInArray(ap.ChannelCode, ccArr) { //渠道码
			return false
		}
	}
	if len(ap.Resolution) > 0 && black.Resolution != "" {
		rlArr := strings.Split(black.Resolution, ",")
		if util.StrInArray(ap.Resolution, rlArr) { //分辨率
			return false
		}
	}
	if len(ap.Language) > 0 && black.Language != "" {
		laArr := strings.Split(black.Language, ",")
		if util.StrInArray(ap.Language, laArr) { //语言
			return false
		}
	}
	if len(ap.StorageSpace) > 0 && black.StorageSpace != "" {
		ssArr := strings.Split(black.StorageSpace, ",")
		if util.StrInArray(ap.StorageSpace, ssArr) { //存储空间
			return false
		}
	}
	if len(ap.AndroidVersionName) > 0 && black.AndroidVersionName != "" {
		avnArr := strings.Split(black.AndroidVersionName, ",")
		if util.StrInArray(ap.AndroidVersionName, avnArr) { //Android 名称
			return false
		}
	}
	if len(ap.AndroidVersionCode) > 0 && black.AndroidVersionCode != "" {
		avcArr := strings.Split(black.AndroidVersionCode, ",")
		if util.StrInArray(ap.AndroidVersionCode, avcArr) { //Android Code
			return false
		}
	}
	if len(ap.SdkName) > 0 && black.SdkName != "" {
		snArr := strings.Split(black.SdkName, ",")
		if util.StrInArray(ap.SdkName, snArr) { //Sdk 名称
			return false
		}
	}
	if len(ap.SdkVersionName) > 0 && black.SdkVersionName != "" {
		svnArr := strings.Split(black.SdkVersionName, ",")
		if util.StrInArray(ap.SdkVersionName, svnArr) { //Sdk 版本名称
			return false
		}
	}
	if len(ap.SdkVersionCode) > 0 && black.SdkVersionCode != "" {
		svcArr := strings.Split(black.SdkVersionCode, ",")
		if util.StrInArray(ap.SdkVersionCode, svcArr) { //Sdk 版本 Code
			return false
		}
	}
	if len(ap.AppName) > 0 && black.AppName != "" {
		anArr := strings.Split(black.AppName, ",")
		if util.StrInArray(ap.AppName, anArr) { //App name
			return false
		}
	}
	if len(ap.AppVersionName) > 0 && black.AppVersionName != "" {
		avnArr := strings.Split(black.AppVersionName, ",")
		if util.StrInArray(ap.AppVersionName, avnArr) { //App 版本名称
			return false
		}
	}
	if len(ap.AppVersionCode) > 0 && black.AppVersionCode != "" {
		avcArr := strings.Split(black.AppVersionCode, ",")
		if util.StrInArray(ap.AppVersionCode, avcArr) { //App 版本 Code
			return false
		}
	}
	if len(ap.CpuArch) > 0 && black.CpuArch != "" {
		caArr := strings.Split(black.CpuArch, ",")
		if util.StrInArray(ap.CpuArch, caArr) { //CPU 架构
			return false
		}
	}

	return true
}

// 验证白名单 - 无黑名单 & 有白名单 - 若在白名单内，则更新；否则不更新；
func checkWhite(ap AppParam, white *AppBlackList) bool {
	if len(ap.DeviceNo) > 0 && white.DeviceNo != "" {
		dnArr := strings.Split(white.DeviceNo, ",")
		if !util.StrInArray(ap.DeviceNo, dnArr) { //  设备号
			return false
		}
	}
	if len(ap.DeviceModel) > 0 && white.DeviceModel != "" {
		dmArr := strings.Split(white.DeviceModel, ",")
		if !util.StrInArray(ap.DeviceModel, dmArr) { //设备型号
			return false
		}
	}
	if len(ap.ChannelCode) > 0 && white.ChannelCode != "" {
		ccArr := strings.Split(white.ChannelCode, ",")
		if !util.StrInArray(ap.ChannelCode, ccArr) { //渠道码
			return false
		}
	}
	if len(ap.Resolution) > 0 && white.Resolution != "" {
		rlArr := strings.Split(white.Resolution, ",")
		if !util.StrInArray(ap.Resolution, rlArr) { //分辨率
			return false
		}
	}
	if len(ap.Language) > 0 && white.Language != "" {
		laArr := strings.Split(white.Language, ",")
		if !util.StrInArray(ap.Language, laArr) { //语言
			return false
		}
	}
	if len(ap.StorageSpace) > 0 && white.StorageSpace != "" {
		ssArr := strings.Split(white.StorageSpace, ",")
		if !util.StrInArray(ap.StorageSpace, ssArr) { //存储空间
			return false
		}
	}
	if len(ap.AndroidVersionName) > 0 && white.AndroidVersionName != "" {
		avnArr := strings.Split(white.AndroidVersionName, ",")
		if !util.StrInArray(ap.AndroidVersionName, avnArr) { //Android 名称
			return false
		}
	}
	if len(ap.AndroidVersionCode) > 0 && white.AndroidVersionCode != "" {
		avcArr := strings.Split(white.AndroidVersionCode, ",")
		if !util.StrInArray(ap.AndroidVersionCode, avcArr) { //Android Code
			return false
		}
	}
	if len(ap.SdkName) > 0 && white.SdkName != "" {
		snArr := strings.Split(white.SdkName, ",")
		if !util.StrInArray(ap.SdkName, snArr) { //Sdk 名称
			return false
		}
	}
	if len(ap.SdkVersionName) > 0 && white.SdkVersionName != "" {
		svnArr := strings.Split(white.SdkVersionName, ",")
		if !util.StrInArray(ap.SdkVersionName, svnArr) { //Sdk 版本名称
			return false
		}
	}
	if len(ap.SdkVersionCode) > 0 && white.SdkVersionCode != "" {
		svcArr := strings.Split(white.SdkVersionCode, ",")
		if !util.StrInArray(ap.SdkVersionCode, svcArr) { //Sdk 版本 Code
			return false
		}
	}
	if len(ap.AppName) > 0 && white.AppName != "" {
		anArr := strings.Split(white.AppName, ",")
		if !util.StrInArray(ap.AppName, anArr) { //App name
			return false
		}
	}
	if len(ap.AppVersionName) > 0 && white.AppVersionName != "" {
		avnArr := strings.Split(white.AppVersionName, ",")
		if !util.StrInArray(ap.AppVersionName, avnArr) { //App 版本名称
			return false
		}
	}
	if len(ap.AppVersionCode) > 0 && white.AppVersionCode != "" {
		avcArr := strings.Split(white.AppVersionCode, ",")
		if !util.StrInArray(ap.AppVersionCode, avcArr) { //App 版本 Code
			return false
		}
	}
	if len(ap.CpuArch) > 0 && white.CpuArch != "" {
		caArr := strings.Split(white.CpuArch, ",")
		if !util.StrInArray(ap.CpuArch, caArr) { //CPU 架构
			return false
		}
	}

	return true
}

// 验证黑&白名单 - 有黑名单 & 有白名单 - 若在白名单内 && 不在黑名单内，则更新，否则不更新；
func checkBlackWhite(ap AppParam, black *AppBlackList, white *AppBlackList) bool {
	if len(ap.DeviceNo) > 0 && white.DeviceNo != "" {
		bdnArr := strings.Split(black.DeviceNo, ",")
		wdnArr := strings.Split(white.DeviceNo, ",")
		if !util.StrInArray(ap.DeviceNo, wdnArr) || util.StrInArray(ap.DeviceNo, bdnArr) { //  设备号
			return false
		}
	}
	if len(ap.DeviceModel) > 0 && white.DeviceModel != "" {
		bdmArr := strings.Split(black.DeviceModel, ",")
		wdmArr := strings.Split(white.DeviceModel, ",")
		if !util.StrInArray(ap.DeviceModel, wdmArr) || util.StrInArray(ap.DeviceModel, bdmArr) { //设备型号
			return false
		}
	}
	if len(ap.ChannelCode) > 0 && white.ChannelCode != "" {
		bccArr := strings.Split(black.ChannelCode, ",")
		wccArr := strings.Split(white.ChannelCode, ",")
		if !util.StrInArray(ap.ChannelCode, wccArr) || util.StrInArray(ap.ChannelCode, bccArr) { //渠道码
			return false
		}
	}
	if len(ap.Resolution) > 0 && white.Resolution != "" {
		brlArr := strings.Split(black.Resolution, ",")
		wrlArr := strings.Split(white.Resolution, ",")
		if !util.StrInArray(ap.Resolution, wrlArr) || util.StrInArray(ap.Resolution, brlArr) { //分辨率
			return false
		}
	}
	if len(ap.Language) > 0 && white.Language != "" {
		blaArr := strings.Split(black.Language, ",")
		wlaArr := strings.Split(white.Language, ",")
		if !util.StrInArray(ap.Language, wlaArr) || util.StrInArray(ap.Language, blaArr) { //语言
			return false
		}
	}
	if len(ap.StorageSpace) > 0 && white.StorageSpace != "" {
		bssArr := strings.Split(black.StorageSpace, ",")
		wssArr := strings.Split(white.StorageSpace, ",")
		if !util.StrInArray(ap.StorageSpace, wssArr) || util.StrInArray(ap.StorageSpace, bssArr) { //存储空间
			return false
		}
	}
	if len(ap.AndroidVersionName) > 0 && white.AndroidVersionName != "" {
		bavnArr := strings.Split(black.AndroidVersionName, ",")
		wavnArr := strings.Split(white.AndroidVersionName, ",")
		if !util.StrInArray(ap.AndroidVersionName, wavnArr) || util.StrInArray(ap.AndroidVersionName, bavnArr) { //Android 名称
			return false
		}
	}
	if len(ap.AndroidVersionCode) > 0 && white.AndroidVersionCode != "" {
		bavcArr := strings.Split(black.AndroidVersionCode, ",")
		wavcArr := strings.Split(white.AndroidVersionCode, ",")
		if !util.StrInArray(ap.AndroidVersionCode, wavcArr) || util.StrInArray(ap.AndroidVersionCode, bavcArr) { //Android Code
			return false
		}
	}
	if len(ap.SdkName) > 0 && white.SdkName != "" {
		bsnArr := strings.Split(black.SdkName, ",")
		wsnArr := strings.Split(white.SdkName, ",")
		if !util.StrInArray(ap.SdkName, wsnArr) || util.StrInArray(ap.SdkName, bsnArr) { //Sdk 名称
			return false
		}
	}
	if len(ap.SdkVersionName) > 0 && white.SdkVersionName != "" {
		bsvnArr := strings.Split(black.SdkVersionName, ",")
		wsvnArr := strings.Split(white.SdkVersionName, ",")
		if !util.StrInArray(ap.SdkVersionName, wsvnArr) || util.StrInArray(ap.SdkVersionName, bsvnArr) { //Sdk 版本名称
			return false
		}
	}
	if len(ap.SdkVersionCode) > 0 && white.SdkVersionCode != "" {
		bsvcArr := strings.Split(black.SdkVersionCode, ",")
		wsvcArr := strings.Split(white.SdkVersionCode, ",")
		if !util.StrInArray(ap.SdkVersionCode, wsvcArr) || util.StrInArray(ap.SdkVersionCode, bsvcArr) { //Sdk 版本 Code
			return false
		}
	}
	if len(ap.AppName) > 0 && white.AppName != "" {
		banArr := strings.Split(black.AppName, ",")
		wanArr := strings.Split(white.AppName, ",")
		if !util.StrInArray(ap.AppName, wanArr) || util.StrInArray(ap.AppName, banArr) { //App name
			return false
		}
	}
	if len(ap.AppVersionName) > 0 && white.AppVersionName != "" {
		bavnArr := strings.Split(black.AppVersionName, ",")
		wavnArr := strings.Split(white.AppVersionName, ",")
		if !util.StrInArray(ap.AppVersionName, wavnArr) || util.StrInArray(ap.AppVersionName, bavnArr) { //App 版本名称
			return false
		}
	}
	if len(ap.AppVersionCode) > 0 && white.AppVersionCode != "" {
		bavcArr := strings.Split(black.AppVersionCode, ",")
		wavcArr := strings.Split(white.AppVersionCode, ",")
		if !util.StrInArray(ap.AppVersionCode, wavcArr) || util.StrInArray(ap.AppVersionCode, bavcArr) { //App 版本 Code
			return false
		}
	}
	if len(ap.CpuArch) > 0 && white.CpuArch != "" {
		bcaArr := strings.Split(black.CpuArch, ",")
		wcaArr := strings.Split(white.CpuArch, ",")
		if !util.StrInArray(ap.CpuArch, wcaArr) || util.StrInArray(ap.CpuArch, bcaArr) { //CPU 架构
			return false
		}
	}

	return true
}
