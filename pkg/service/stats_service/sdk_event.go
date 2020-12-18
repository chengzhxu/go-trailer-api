package stats_service

import (
	"encoding/json"
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/util"
	"strings"
)

type SdkEvent struct {
	ClientTime string `json:"client_time" binding:"nm_bas_time"` //客户端时间 格式：2020-12-12 12:12:12
	DeviceNo   string `json:"device_no" binding:"required"`      //设备号
	IMEI       string `json:"imei" binding:""`                   //IMEI
	//IDFA 						string `json:"idfa" `
	NewpUid        string `json:"newpuid" `                                        //APP 用户账号 ID
	NewSessionId   string `json:"newsession_id" binding:"required"`                //会话 ID
	ScreenWidth    int    `json:"screen_width" binding:"required"`                 //分辨率 宽度
	ScreenHeight   int    `json:"screen_height" binding:"required"`                //分辨率 高度
	OsVersionName  string `json:"os_version_name" binding:"required"`              //操作系统版本名称
	OsVersionCode  string `json:"os_version_code" binding:"required"`              //操作系统版本 Code
	DeviceBrand    string `json:"device_brand" `                                   //设备品牌
	DeviceModel    string `json:"device_model" binding:"required"`                 //设备型号
	ChannelCode    string `json:"channel_code" binding:"required"`                 //渠道码
	AppName        string `json:"app_name" binding:"required"`                     //APP 名称
	AppVersionCode string `json:"app_version_code" binding:"required"`             //APP 版本 Code
	AppVersionName string `json:"app_version_name" binding:"required"`             //APP 版本名称
	SdkName        string `json:"sdk_name" binding:"required"`                     //SDK 名称
	SdkVersionName string `json:"sdk_version_name" binding:"required"`             //SDK 版本名称
	SdkVersionCode string `json:"sdk_version_code" binding:"required"`             //SDK 版本 Code
	PageName       string `json:"page_name" `                                      //所在页面
	IP             string `json:"ip" `                                             //IP
	NetType        string `json:"net_type" binding:"required"`                     //网络类型	WIFI/4G/5G
	NewEventType   *int   `json:"newevent_type" binding:"required,sdk_event_type"` //事件类型  0:自定义事件,1:预置事件
	EventName      string `json:"event_name" binding:"required"`                   //事件名称
	EventKvJson    string `json:"event_kv_json" binding:"required,sdk_event_kt"`   //参数和参数值数据 -json数组 - 格式：{"key1": "val1", "key2": "val2"}
}

type ObjSdkEvents struct {
	SdkEvents string `json:"sdk_events" binding:"required"`
}

type ArrSdkEvents struct {
	SdkEvents []SdkEvent
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
		se.ClientTime = util.GetCurrentTime() //服务器时间
	}
	sdkEvent := mapSdkEvent(se)

	if err := model.InsertSdkEvent(sdkEvent); err != nil {
		return err
	}

	return nil
}

//解析事件参数，拆分为上报事件
func TranceEventKtJson(se SdkEvent) []SdkEvent {
	eventArr := []SdkEvent{}

	event_kv_json := se.EventKvJson
	if event_kv_json != "" {
		ej := strings.NewReader(event_kv_json)
		dec := json.NewDecoder(ej)
		switch event_kv_json[0] {
		case 91:
			// "[" 开头的Json（数组型Json）
			param := []interface{}{}
			dec.Decode(&param)
			for _, v := range param {
				eventJsonStr, _ := json.Marshal(v)
				se.EventKvJson = string(eventJsonStr)
				eventArr = append(eventArr, se)
			}

			return eventArr
		case 123:
			// "{" 开头的Json（对象型Json）
		}
	}
	eventArr = append(eventArr, se)

	return eventArr
}
