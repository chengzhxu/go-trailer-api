package app_service

import "go-trailer-api/pkg/model"

type AppParam struct {
	DeviceNo           string `json:"device_no" binding:"required"`            //设备号
	DeviceModel        string `json:"device_model" binding:"required"`         //设备型号
	ChannelCode        string `json:"channel_code" binding:"required"`         //渠道码
	Resolution         string `json:"resolution" binding:""`                   //分辨率
	Language           string `json:"language" binding:""`                     //语言
	StorageSpace       string `json:"storage_space" binding:"required"`        //存储空间
	AndroidVersionName string `json:"android_version_name" binding:"required"` //Android 版本名称
	AndroidVersionCode string `json:"android_version_code" binding:"required"` //Android 版本 Code
	SdkName            string `json:"sdk_name" binding:"required"`             //SDK 名称
	SdkVersionName     string `json:"sdk_version_name" binding:"required"`     //SDK 版本名称
	SdkVersionCode     string `json:"sdk_version_code" binding:""`             //SDK 版本 Code
	AppName            string `json:"app_name" binding:""`                     //APP 名称
	AppVersionName     string `json:"app_version_name" binding:""`             //APP 版本名称
	AppVersionCode     string `json:"app_version_code" binding:""`             //APP 版本 Code
	CpuArch            string `json:"cpu_arch" binding:"required"`             //CPU 架构
	IsHotUpdate        int    `json:"is_hot_update" binding:""`                //是否热更  0:否  1:是
	IsSecure           bool   `json:"isSecure" binding:""`                     //判断返回链接形式 https or http;
}

func mapAppParam(ap *AppParam) map[string]interface{} {
	return map[string]interface{}{
		"device_no":            ap.DeviceNo,
		"device_model":         ap.DeviceModel,
		"channel_code":         ap.ChannelCode,
		"resolution":           ap.Resolution,
		"language":             ap.Language,
		"storage_space":        ap.StorageSpace,
		"android_version_name": ap.AndroidVersionName,
		"android_version_code": ap.AndroidVersionCode,
		"sdk_name":             ap.SdkName,
		"sdk_version_code":     ap.SdkVersionCode,
		"sdk_version_name":     ap.SdkVersionName,
		"app_name":             ap.AppName,
		"app_version_code":     ap.AppVersionCode,
		"app_version_name":     ap.AppVersionName,
		"cpu_arch":             ap.CpuArch,
		"is_hot_update":        ap.IsHotUpdate,
		"is_secure":            ap.IsSecure,
	}
}

// 获取 App 最新版本
func (ap *AppParam) GetNewAppVersion() (*model.AppVersion, error) {
	appParam := mapAppParam(ap)
	appVersion, err := model.GetNewAppVersion(appParam)
	if err != nil {
		return nil, err
	}

	return appVersion, nil
}
