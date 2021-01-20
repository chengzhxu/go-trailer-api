package model

import "go-trailer-api/pkg/util"

type AppLogs struct {
	Id            int    `json:"id" gorm:"column:id"`
	URL           string `json:"url" gorm:"column:url"`
	DeviceNo      string `json:"device_no" gorm:"column:device_no"`
	OsVersionCode string `json:"os_version_code" gorm:"column:os_version_code"`
	ChannelCode   string `json:"channel_code" gorm:"column:channel_code"`
	CreateTime    string `json:"create_time" gorm:"column:create_time"`
}

func (AppLogs) TableName() string {
	return "stats_app_logs"
}

func InsertAppLog(data map[string]interface{}) error {
	// todo mapTo function
	AppLog := AppLogs{
		URL:           data["url"].(string),
		DeviceNo:      data["device_no"].(string),
		OsVersionCode: data["os_version_code"].(string),
		ChannelCode:   data["channel_code"].(string),
		CreateTime:    util.GetCurrentTime(),
	}
	if err := db.Create(&AppLog).Error; err != nil {
		return err
	}

	return nil
}
