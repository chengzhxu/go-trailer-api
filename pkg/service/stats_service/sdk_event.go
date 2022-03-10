package stats_service

import (
	"encoding/json"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/util"
	"io"
	"os"
)

type SdkEvent struct {
	ClientTime string `json:"client_time" binding:"nm_bas_time"` //客户端时间 格式：2020-12-12 12:12:12
	DeviceNo   string `json:"device_no" binding:"required"`      //设备号
	IMEI       string `json:"imei" binding:""`                   //IMEI
	//IDFA 						string `json:"idfa" `
	NewpUid        string `json:"newpuid" `                            //APP 用户账号 ID
	NewSessionId   string `json:"newsession_id" binding:"required"`    //会话 ID
	ScreenWidth    int    `json:"screen_width" binding:"required"`     //分辨率 宽度
	ScreenHeight   int    `json:"screen_height" binding:"required"`    //分辨率 高度
	OsVersionName  string `json:"os_version_name" binding:"required"`  //操作系统版本名称
	OsVersionCode  string `json:"os_version_code" binding:"required"`  //操作系统版本 Code
	DeviceBrand    string `json:"device_brand" `                       //设备品牌
	DeviceModel    string `json:"device_model" binding:"required"`     //设备型号
	ChannelCode    string `json:"channel_code" binding:"required"`     //渠道码
	AppName        string `json:"app_name" binding:"required"`         //APP 名称
	AppVersionCode string `json:"app_version_code" binding:"required"` //APP 版本 Code
	AppVersionName string `json:"app_version_name" binding:"required"` //APP 版本名称
	SdkName        string `json:"sdk_name" binding:"required"`         //SDK 名称
	SdkVersionName string `json:"sdk_version_name" binding:"required"` //SDK 版本名称
	SdkVersionCode string `json:"sdk_version_code" binding:"required"` //SDK 版本 Code
	PageName       string `json:"page_name" `                          //所在页面
	IP             string `json:"ip" `                                 //IP
	MAC            string `json:"mac" `                                //MAC 地址
	CpuArch        string `json:"cpu_arch" `                           //CpuArch
	NetType        string `json:"net_type" binding:"required"`         //网络类型	WIFI/4G/5G
	NewEventType   int    `json:"newevent_type" binding:""`            //事件类型  0:自定义事件,1:预置事件
	EventName      string `json:"event_name" binding:""`               //事件名称
	EventKvJson    string `json:"event_kv_json" binding:""`            //参数和参数值数据 -json数组 - 格式：{"key1": "val1", "key2": "val2"}
	EventInfo      string `json:"event_info" binding:"sdk_event_info"` //参数和参数值数据 -json数组 - 格式：[{"event_kv_json":{"key1": "val1", "key2": "val2"},"event_name":"sss","event_type":1},...]
	Signature      string `json:"signature" binding:""`                //签名
}

type ObjSdkEvents struct {
	SdkEvents string `json:"sdk_events" binding:"required"`
}

type ArrSdkEvents struct {
	SdkEvents []SdkEvent
}

type EventInfo struct {
	EventName   string      `json:"event_name" binding:"required"`
	EventKvJson interface{} `json:"event_kv_json" binding:"required,sdk_event_kt"`
	EventType   int         `json:"event_type" binding:"required,sdk_event_type"`
}

var eventFilePath = "storage/sdk_events/"

func MapSdkEvent(se SdkEvent) map[string]interface{} {
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
		"mac":              se.MAC,
		"cpu_arch":         se.CpuArch,
		"net_type":         se.NetType,
		"newevent_type":    se.NewEventType,
		"event_name":       se.EventName,
		"event_kv_json":    se.EventKvJson,
		"event_info":       se.EventInfo,
	}
}

func (se SdkEvent) Insert() error {
	if len(se.ClientTime) == 0 {
		se.ClientTime = util.GetCurrentTime() //服务器时间
	}
	sdkEvent := MapSdkEvent(se)

	// 写入 MySql
	if err := model.InsertSdkEvent(sdkEvent); err != nil {
		return err
	}
	// 写入文件，上传至 ES
	if err := sdkEventWriteFile(sdkEvent); err != nil {
		return err
	}

	return nil
}

//解析事件参数，拆分为上报事件
func TranceEventKtJson(se SdkEvent) []SdkEvent {
	eventArr := []SdkEvent{}

	eventInfo := se.EventInfo
	if eventInfo != "" {
		var arrEventInfo []EventInfo
		err := json.Unmarshal([]byte(eventInfo), &arrEventInfo)
		if err != nil {
			logging.Error(err)
		} else {
			for _, v := range arrEventInfo {
				eventJsonStr, _ := json.Marshal(v.EventKvJson)
				se.EventKvJson = string(eventJsonStr)
				se.EventName = v.EventName
				se.NewEventType = v.EventType
				eventArr = append(eventArr, se)
			}
		}
	} else {
		eventArr = append(eventArr, se)
	}

	return eventArr
}

// 事件写入文件
func sdkEventWriteFile(se map[string]interface{}) error {
	if se["device_no"] != "" && se["channel_code"] != "" {
		jsonStr, err := json.Marshal(se)
		if err != nil {
			return err
		}

		_, pathErr := os.Stat(eventFilePath)
		if pathErr != nil {
			os.Mkdir(eventFilePath, os.ModePerm)
		}

		fileName := util.GetCurrentDate() + ".log"
		filePath := eventFilePath + fileName
		var file *os.File

		_, fileErr := os.Stat(filePath)
		if fileErr != nil {
			file, _ = os.Create(filePath)
		} else {
			file, _ = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			defer file.Close()
		}

		_, writeErr := io.WriteString(file, string(jsonStr)+"\r\n")
		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}
