package model

import (
	"go-trailer-api/pkg/tool"
)

type SdkErrors struct {
	Id             int    `json:"id" gorm:"column:id"`
	DeviceNo       string `json:"device_no" gorm:"column:device_no"`
	ChannelCode    string `json:"channel_code" gorm:"column:channel_code"`
	CrashTime      string `json:"crash_time" gorm:"column:crash_time"`
	CrashLog       string `json:"crash_log" gorm:"column:crash_log"`
	SdkName        string `json:"sdk_name" gorm:"column:sdk_name"`
	SdkVersionName string `json:"sdk_version_name" gorm:"column:sdk_version_name"`
	SdkVersionCode string `json:"sdk_version_code" gorm:"column:sdk_version_code"`
	AppName        string `json:"app_name" gorm:"column:app_name"`
	AppVersionName string `json:"app_version_name" gorm:"column:app_version_name"`
	AppVersionCode string `json:"app_version_code" gorm:"column:app_version_code"`
	UserId         int    `json:"user_id" gorm:"column:user_id"`
	Ext            string `json:"ext" gorm:"column:ext"`
	CreateTime     string `json:"create_time" gorm:"column:create_time"`
}

func (SdkErrors) TableName() string {
	return "stats_sdk_errors"
}

func InsertSdkError(data map[string]interface{}) error {
	// todo mapTo function
	sdkErr := SdkErrors{
		DeviceNo:       data["device_no"].(string),
		ChannelCode:    data["channel_code"].(string),
		CrashTime:      data["crash_time"].(string),
		CrashLog:       data["crash_log"].(string),
		SdkName:        data["sdk_name"].(string),
		SdkVersionName: data["sdk_version_name"].(string),
		SdkVersionCode: data["sdk_version_code"].(string),
		AppName:        data["app_name"].(string),
		AppVersionName: data["app_version_name"].(string),
		AppVersionCode: data["app_version_code"].(string),
		UserId:         data["user_id"].(int),
		Ext:            data["ext"].(string),
		CreateTime:     tool.GetCurrentTime(),
	}
	if err := db.Create(&sdkErr).Error; err != nil {
		return err
	}

	return nil
}
