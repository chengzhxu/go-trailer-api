package stats_service

import "go-trailer-api/pkg/model"

type SdkError struct {
	DeviceNo       string `json:"device_no" binding:"required"`           //设备号
	ChannelCode    string `json:"channel_code" binding:"required"`        //渠道码
	CrashTime      string `json:"crash_time" binding:"required,bas_time"` //Crash 时间
	CrashLog       string `json:"crash_log" binding:"required"`           //Crash 日志
	SdkName        string `json:"sdk_name" binding:"required"`            //SDK 名称
	SdkVersionName string `json:"sdk_version_name" binding:"required"`    //SDK 版本名称
	SdkVersionCode string `json:"sdk_version_code" binding:"required"`    //SDK 版本 Code
	AppName        string `json:"app_name" binding:"required"`            //APP 名称
	AppVersionName string `json:"app_version_name" binding:"required"`    //APP 版本名称
	AppVersionCode string `json:"app_version_code" binding:"required"`    //APP 版本 Code
	UserId         int    `json:"user_id" binding:""`                     //USER ID
	Ext            string `json:"ext" binding:""`                         //自定义数据
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
