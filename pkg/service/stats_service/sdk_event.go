package stats_service

import (
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/tool"
)

type SdkEvent struct {
	ClientTime string `json:"client_time" binding:"nm_bas_time"`
	DeviceNo   string `json:"device_no" binding:"required"`
	IMEI       string `json:"imei" binding:"required"`
	//IDFA 						string `json:"idfa" `
	NewpUid        string `json:"newpuid" `
	NewSessionId   string `json:"newsession_id" binding:"required"`
	ScreenWidth    int    `json:"screen_width" binding:"required"`
	ScreenHeight   int    `json:"screen_height" binding:"required"`
	OsVersionName  string `json:"os_version_name" binding:"required"`
	OsVersionCode  string `json:"os_version_code" binding:"required"`
	DeviceBrand    string `json:"device_brand" `
	DeviceModel    string `json:"device_model" binding:"required"`
	ChannelCode    string `json:"channel_code" binding:"required"`
	AppName        string `json:"app_name" binding:"required"`
	AppVersionCode string `json:"app_version_code" binding:"required"`
	AppVersionName string `json:"app_version_name" binding:"required"`
	SdkName        string `json:"sdk_name" binding:"required"`
	SdkVersionName string `json:"sdk_version_name" binding:"required"`
	SdkVersionCode string `json:"sdk_version_code" binding:"required"`
	PageName       string `json:"page_name" `
	IP             string `json:"ip" `
	NetType        string `json:"net_type" binding:"required"`
	NewEventType   *int   `json:"newevent_type" binding:"required,sdk_event_type"`
	EventName      string `json:"event_name" binding:"required"`
	EventKvJson    string `json:"event_kv_json" binding:"required,sdk_event_kt"`
}

func mapSdkEvent(se *SdkEvent) map[string]interface{} {
	return map[string]interface{}{
		"client_time":      se.ClientTime,
		"device_no":        se.DeviceNo,
		"imei":             se.IMEI,
		"newpuid":          se.NewpUid,
		"newsession_id":    se.NewSessionId,
		"screen_width":     se.ScreenWidth,
		"screen_height":    se.ScreenHeight,
		"os_version_name":  se.OsVersionName,
		"os_version_code":  se.OsVersionCode,
		"device_brand":     se.DeviceBrand,
		"device_model":     se.DeviceModel,
		"channel_code":     se.ChannelCode,
		"app_name":         se.AppName,
		"app_version_code": se.AppVersionCode,
		"app_version_name": se.AppVersionName,
		"sdk_name":         se.SdkName,
		"sdk_version_name": se.SdkVersionName,
		"sdk_version_code": se.SdkVersionCode,
		"page_name":        se.PageName,
		"ip":               se.IP,
		"net_type":         se.NetType,
		"newevent_type":    *se.NewEventType,
		"event_name":       se.EventName,
		"event_kv_json":    se.EventKvJson,
	}
}

func (se *SdkEvent) Insert() error {
	if len(se.ClientTime) == 0 {
		se.ClientTime = tool.GetCurrentTime() //服务器时间
	}
	sdkEvent := mapSdkEvent(se)

	if err := model.InsertSdkEvent(sdkEvent); err != nil {
		return err
	}

	return nil
}
