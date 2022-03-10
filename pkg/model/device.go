package model

import (
	"go-trailer-api/pkg/util"
)

type Device struct {
	Id                 int    `json:"id" gorm:"column:id"`
	DeviceNo           string `json:"device_no" gorm:"column:device_no"`
	DeviceVendor       string `json:"device_vendor" gorm:"column:device_vendor"`
	DeviceModel        string `json:"device_model" gorm:"column:device_model"`
	DeviceName         string `json:"device_name" gorm:"column:device_name"`
	ChannelCode        string `json:"channel_code" gorm:"column:channel_code"`
	Resolution         string `json:"resolution" gorm:"column:resolution"`
	AndroidVersionCode string `json:"android_version_code" gorm:"column:android_version_code"`
	AndroidVersionName string `json:"android_version_name" gorm:"column:android_version_name"`
	AppName            string `json:"app_name" gorm:"column:app_name"`
	AppVersionCode     string `json:"app_version_code" gorm:"column:app_version_code"`
	AppVersionName     string `json:"app_version_name" gorm:"column:app_version_name"`
	IP                 string `json:"ip" gorm:"column:ip"`
	MAC                string `json:"mac" gorm:"column:mac"`
	CpuArch            string `json:"cpu_arch" gorm:"column:cpu_arch"`
	CreateTime         string `json:"create_time" gorm:"column:create_time"`
}

func (Device) TableName() string {
	return "device"
}

func InsertDevice(data map[string]interface{}) error {
	// todo mapTo function
	device := Device{
		DeviceNo:           data["device_no"].(string),
		DeviceVendor:       data["device_vendor"].(string),
		DeviceModel:        data["device_model"].(string),
		DeviceName:         data["device_name"].(string),
		ChannelCode:        data["channel_code"].(string),
		Resolution:         data["resolution"].(string),
		AndroidVersionName: data["android_version_name"].(string),
		AndroidVersionCode: data["android_version_code"].(string),
		AppName:            data["app_name"].(string),
		AppVersionName:     data["app_version_name"].(string),
		AppVersionCode:     data["app_version_code"].(string),
		IP:                 data["ip"].(string),
		MAC:                data["mac"].(string),
		CpuArch:            data["cpu_arch"].(string),
		CreateTime:         util.GetCurrentTime(),
	}
	if err := db.Create(&device).Error; err != nil {
		return err
	}

	return nil
}
