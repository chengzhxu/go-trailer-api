package stats_service

import "go-trailer-api/pkg/model"

type Device struct {
	DeviceNo           string `json:"device_no" binding:"required"`            //设备号
	DeviceVendor       string `json:"device_vendor" binding:"required"`        //设备厂商
	DeviceModel        string `json:"device_model" binding:"required"`         //设备型号
	DeviceName         string `json:"device_name" binding:"required"`          //设备名称
	ChannelCode        string `json:"channel_code" binding:"required"`         //渠道码
	Resolution         string `json:"resolution" binding:"required"`           //分辨率
	AndroidVersionName string `json:"android_version_name" binding:"required"` //Android 版本名称
	AndroidVersionCode string `json:"android_version_code" binding:"required"` //Android 版本 Code
	AppName            string `json:"app_name" binding:"required"`             //APP 名称
	AppVersionCode     string `json:"app_version_code" binding:"required"`     //APP 版本 Code
	AppVersionName     string `json:"app_version_name" binding:"required"`     //APP 版本名称
	IP                 string `json:"ip" binding:"required"`                   //IP
	MAC                string `json:"mac" binding:""`                          //MAC
	Signature          string `json:"signature" binding:""`                    //签名
}

func MapDevice(d Device) map[string]interface{} {
	return map[string]interface{}{
		"device_no":            d.DeviceNo,
		"device_vendor":        d.DeviceVendor,
		"device_model":         d.DeviceModel,
		"device_name":          d.DeviceName,
		"channel_code":         d.ChannelCode,
		"resolution":           d.Resolution,
		"android_version_name": d.AndroidVersionName,
		"android_version_code": d.AndroidVersionCode,
		"app_name":             d.AppName,
		"app_version_code":     d.AppVersionCode,
		"app_version_name":     d.AppVersionName,
		"ip":                   d.IP,
		"mac":                  d.MAC,
	}
}

func (d Device) Insert() error {
	device := MapDevice(d)

	if err := model.InsertDevice(device); err != nil {
		return err
	}

	return nil
}
