package model

//App 版本更新黑/白名单规则
type AppBlackList struct {
	IP                 string `json:"ip" gorm:"column:ip"`                                     //IP
	DeviceNo           string `json:"device_no" gorm:"column:device_no"`                       //设备编号
	DeviceModel        string `json:"device_model" gorm:"column:device_model"`                 //设备型号
	ChannelCode        string `json:"channel_code" gorm:"column:channel_code"`                 //渠道码
	Resolution         string `json:"resolution" gorm:"column:resolution"`                     //分辨率
	Language           string `json:"language" gorm:"column:language"`                         //语言
	StorageSpace       string `json:"storage_space" gorm:"column:storage_space"`               //存储空间
	AndroidVersionCode string `json:"android_version_code" gorm:"column:android_version_code"` //Android 版本 Code
	AndroidVersionName string `json:"android_version_name" gorm:"column:android_version_name"` //Android 版本 名称
	AppName            string `json:"app_name" gorm:"column:app_name"`                         //App 名称
	AppVersionCode     string `json:"app_version_code" gorm:"column:app_version_code"`         //App 版本 Code
	AppVersionName     string `json:"app_version_name" gorm:"column:app_version_name"`         //App 版本 名称
	SdkName            string `json:"sdk_name" gorm:"column:sdk_name"`                         //Sdk 名称
	SdkVersionName     string `json:"sdk_version_name" gorm:"column:sdk_version_name"`         //Sdk 版本名称
	SdkVersionCode     string `json:"sdk_version_code" gorm:"column:sdk_version_code"`         //Sdk 版本 Code
	CpuArch            string `json:"cpu_arch" gorm:"column:cpu_arch"`                         //Cpu 架构
}
