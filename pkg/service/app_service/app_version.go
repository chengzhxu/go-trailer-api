package app_service

import "go-trailer-api/pkg/model"

type AppParam struct {
	DeviceNo           string `json:"device_no" binding:"required"`
	DeviceModel        string `json:"device_model" binding:"required"`
	ChannelCode        string `json:"channel_code" binding:"required"`
	Resolution         string `json:"resolution" binding:""`
	Language           string `json:"language" binding:""`
	StorageSpace       string `json:"storage_space" binding:"required"`
	AndroidVersionName string `json:"android_version_name" binding:"required"`
	AndroidVersionCode string `json:"android_version_code" binding:"required"`
	SdkName            string `json:"sdk_name" binding:"required"`
	SdkVersionName     string `json:"sdk_version_name" binding:"required"`
	SdkVersionCode     string `json:"sdk_version_code" binding:""`
	AppName            string `json:"app_name" binding:""`
	AppVersionName     string `json:"app_version_name" binding:""`
	AppVersionCode     string `json:"app_version_code" binding:""`
	CpuArch            string `json:"cpu_arch" binding:"required"`
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
