package stats_service

import "go-trailer-api/pkg/model"

type Device struct {
	DeviceNo           string `json:"device_no" binding:"required"`
	DeviceVendor       string `json:"device_vendor" binding:"required"`
	DeviceModel        string `json:"device_model" binding:"required"`
	DeviceName         string `json:"device_name" binding:"required"`
	ChannelCode        string `json:"channel_code" binding:"required"`
	Resolution         string `json:"resolution" binding:"required"`
	AndroidVersionName string `json:"android_version_name" binding:"required"`
	AndroidVersionCode string `json:"android_version_code" binding:"required"`
	AppName            string `json:"app_name" binding:"required"`
	AppVersionCode     string `json:"app_version_code" binding:"required"`
	AppVersionName     string `json:"app_version_name" binding:"required"`
	IP                 string `json:"ip" binding:"required"`
}

func mapDevice(d *Device) map[string]interface{} {
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
	}
}

func (d *Device) Insert() error {
	device := mapDevice(d)

	if err := model.InsertDevice(device); err != nil {
		return err
	}

	return nil
}
