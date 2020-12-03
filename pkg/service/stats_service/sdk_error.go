package stats_service

import "go-trailer-api/pkg/model"

type SdkError struct {
	DeviceNo       string `json:"device_no" binding:"required"`
	ChannelCode    string `json:"channel_code" binding:"required"`
	CrashTime      string `json:"crash_time" binding:"required,bas_time"`
	CrashLog       string `json:"crash_log" binding:"required"`
	SdkName        string `json:"sdk_name" binding:"required"`
	SdkVersionName string `json:"sdk_version_name" binding:"required"`
	SdkVersionCode string `json:"sdk_version_code" binding:"required"`
	AppName        string `json:"app_name" binding:"required"`
	AppVersionName string `json:"app_version_name" binding:"required"`
	AppVersionCode string `json:"app_version_code" binding:"required"`
	UserId         int    `json:"user_id" binding:""`
	Ext            string `json:"ext" binding:""`
}

func mapSdkError(se *SdkError) map[string]interface{} {
	return map[string]interface{}{
		"device_no":        se.DeviceNo,
		"channel_code":     se.ChannelCode,
		"crash_time":       se.CrashTime,
		"crash_log":        se.CrashLog,
		"sdk_name":         se.SdkName,
		"sdk_version_name": se.SdkVersionName,
		"sdk_version_code": se.SdkVersionCode,
		"app_name":         se.AppName,
		"app_version_code": se.AppVersionCode,
		"app_version_name": se.AppVersionName,
		"user_id":          se.UserId,
		"ext":              se.Ext,
	}
}

func (se *SdkError) Insert() error {
	sdkErr := mapSdkError(se)

	if err := model.InsertSdkError(sdkErr); err != nil {
		return err
	}

	return nil
}
