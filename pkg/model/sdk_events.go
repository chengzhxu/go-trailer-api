package model

import (
	"encoding/json"
	"go-trailer-api/pkg/gredis"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/util"
	"strconv"
)

type SdkEvents struct {
	Id             int    `json:"id" gorm:"column:id"`
	ClientTime     string `json:"client_time" gorm:"column:client_time"`
	DeviceNo       string `json:"device_no" gorm:"column:device_no"`
	IMEI           string `json:"imei" gorm:"column:imei"`
	IDFA           string `json:"idfa" gorm:"column:idfa"`
	NewpUid        string `json:"newpuid" gorm:"column:newpuid"`
	NewSessionId   string `json:"newsession_id" gorm:"column:newsession_id"`
	ScreenWidth    int    `json:"screen_width" gorm:"column:screen_width"`
	ScreenHeight   int    `json:"screen_height" gorm:"column:screen_height"`
	OsVersionName  string `json:"os_version_name" gorm:"column:os_version_name"`
	OsVersionCode  string `json:"os_version_code" gorm:"column:os_version_code"`
	DeviceBrand    string `json:"device_brand" gorm:"column:device_brand"`
	DeviceModel    string `json:"device_model" gorm:"column:device_model"`
	ChannelCode    string `json:"channel_code" gorm:"column:channel_code"`
	AppName        string `json:"app_name" gorm:"column:app_name"`
	AppVersionCode string `json:"app_version_code" gorm:"column:app_version_code"`
	AppVersionName string `json:"app_version_name" gorm:"column:app_version_name"`
	SdkName        string `json:"sdk_name" gorm:"column:sdk_name"`
	SdkVersionName string `json:"sdk_version_name" gorm:"column:sdk_version_name"`
	SdkVersionCode string `json:"sdk_version_code" gorm:"column:sdk_version_code"`
	PageName       string `json:"page_name" gorm:"column:page_name"`
	IP             string `json:"ip" gorm:"column:ip"`
	NetType        string `json:"net_type" gorm:"column:net_type"`
	NewEventType   int    `json:"newevent_type" gorm:"column:newevent_type"`
	EventName      string `json:"event_name" gorm:"column:event_name"`
	EventKvJson    string `json:"event_kv_json" gorm:"column:event_kv_json"`
	CreateTime     string `json:"create_time" gorm:"column:create_time"`
}

type EventKv struct {
	TrailerId int `json:"trailer_id" `
}

var TrailerExposureKey = "trailer:exposure" //记录预告片曝光 redis key

func (SdkEvents) TableName() string {
	return "stats_sdk_events"
}

func InsertSdkEvent(data map[string]interface{}) error {
	// todo mapTo function
	event := SdkEvents{
		ClientTime:     data["client_time"].(string),
		DeviceNo:       data["device_no"].(string),
		IMEI:           data["imei"].(string),
		NewpUid:        data["newpuid"].(string),
		NewSessionId:   data["newsession_id"].(string),
		ScreenWidth:    data["screen_width"].(int),
		ScreenHeight:   data["screen_height"].(int),
		OsVersionName:  data["os_version_name"].(string),
		OsVersionCode:  data["os_version_code"].(string),
		DeviceBrand:    data["device_brand"].(string),
		DeviceModel:    data["device_model"].(string),
		ChannelCode:    data["channel_code"].(string),
		AppName:        data["app_name"].(string),
		AppVersionCode: data["app_version_code"].(string),
		AppVersionName: data["app_version_name"].(string),
		SdkName:        data["sdk_name"].(string),
		SdkVersionName: data["sdk_version_name"].(string),
		SdkVersionCode: data["sdk_version_code"].(string),
		PageName:       data["page_name"].(string),
		IP:             data["ip"].(string),
		NetType:        data["net_type"].(string),
		NewEventType:   data["newevent_type"].(int),
		EventName:      data["event_name"].(string),
		EventKvJson:    data["event_kv_json"].(string),
		CreateTime:     util.GetCurrentTime(),
	}

	if err := db.Create(&event).Error; err != nil {
		return err
	}
	RecordTrailerExposureNum(event) //累加预告片曝光数

	return nil
}

//记录预告片 曝光
func RecordTrailerExposureNum(event SdkEvents) {
	if event.EventName == "播放预告片" || event.EventName == "player_trailer" { //播放事件
		if event.EventKvJson != "" {
			e := EventKv{}
			err := json.Unmarshal([]byte(event.EventKvJson), &e)
			if err != nil {
				logging.Error(err.Error())
			} else { //记录曝光 and 播放次数（用于频控）
				trailerId := e.TrailerId
				if trailerId > 0 {
					conn := gredis.RedisConn.Get() //获取 Redis
					trailerIdStr := strconv.Itoa(trailerId)
					key := "trailer-" + trailerIdStr + ":hits"
					_, err = conn.Do("HINCRBY", TrailerExposureKey, key, 1)
					if err != nil {
						logging.Error("累加曝光(" + key + "): " + err.Error())
					}
					if event.DeviceNo != "" { //记录播放次数， 用于频控
						TrailerCapKey := gredis.TrailerCapBeginKey + event.DeviceNo + gredis.TrailerCapEndKey
						CapCountKey := gredis.TrailerCapBeginKey + trailerIdStr + gredis.TrailerCapCountEndKey
						CapDateKey := gredis.TrailerCapBeginKey + trailerIdStr + gredis.TrailerCapDateEndKey
						_, err = conn.Do("HINCRBY", TrailerCapKey, CapCountKey, 1)
						if err != nil {
							logging.Error("累加播放-频控(" + key + "): " + err.Error())
						}
						date, _ := conn.Do("hget", TrailerCapKey, CapDateKey)
						if date == nil { //记录首次播放时间 - 用于清理频控
							_, err = conn.Do("hset", TrailerCapKey, CapDateKey, util.GetCurrentDate())
							if err != nil {
								logging.Error("累加播放-首次播放(" + key + "): " + err.Error())
							}
						}
					}
				}
			}
		}
	}
}
