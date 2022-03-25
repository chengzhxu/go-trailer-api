package stats_service

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/service/aliyun"
	"go-trailer-api/pkg/setting"
	"mime/multipart"
	"strings"
)

type AppLog struct {
	URL           string `form:"url" json:"url" binding:""`                                 //Log url
	DeviceNo      string `form:"device_no" json:"device_no" binding:"required"`             //设备号
	OsVersionCode string `form:"os_version_code" json:"os_version_code" binding:"required"` //系统版本
	ChannelCode   string `form:"channel_code" json:"channel_code" binding:"required"`       //渠道号
	LogType       int    `form:"log_type" json:"log_type" binding:"required,app_log_type"`  //es 类型    1:沙发桌面  2:布丁屏保
}

func MapAppLog(log AppLog) map[string]interface{} {
	return map[string]interface{}{
		"url":             log.URL,
		"device_no":       log.DeviceNo,
		"os_version_code": log.OsVersionCode,
		"channel_code":    log.ChannelCode,
	}
}

func (log AppLog) Insert() error {
	appLog := MapAppLog(log)

	if err := model.InsertAppLog(appLog); err != nil {
		return err
	}

	return nil
}

// 上传日志文件到阿里云 OSS
func UploadLogToAlyOss(header *multipart.FileHeader, logType int, c *gin.Context) (string, error) {
	ossPath := "" //OSS 存储路径

	// 读取文件
	file, err := header.Open()
	if err != nil {
		return ossPath, err
	}

	ossPath, ossErr := aliyun.UploadFileToAlyOss(file, header.Filename, logType)
	if ossErr != nil {
		logging.Error(ossErr)
		return ossPath, ossErr
	}

	return ossPath, nil
}

// 获取 Log 设备 & 渠道 白名单
func GetAppLogWhiteList() map[string]interface{} {
	whiteList := make(map[string]interface{})

	deviceNo := strings.ReplaceAll(setting.AppLogWhiteListSetting.DeviceNo, " ", "")       //设备号
	channelCode := strings.ReplaceAll(setting.AppLogWhiteListSetting.ChannelCode, " ", "") //渠道号

	deviceNoArr := []string{}
	channelCodeArr := []string{}

	if deviceNo != "" {
		deviceNoArr = strings.Split(deviceNo, ",")
	}
	if channelCode != "" {
		channelCodeArr = strings.Split(channelCode, ",")
	}

	whiteList["device_no"] = deviceNoArr
	whiteList["channel_code"] = channelCodeArr

	return whiteList
}
