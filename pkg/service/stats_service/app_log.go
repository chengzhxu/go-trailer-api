package stats_service

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/service/aliyun"
	"mime/multipart"
	"os"
)

type AppLog struct {
	URL           string `json:"url" binding:""`                     //Log url
	DeviceNo      string `json:"device_no" binding:"required"`       //设备号
	OsVersionCode string `json:"os_version_code" binding:"required"` //系统版本
	ChannelCode   string `json:"channel_code" binding:"required"`    //渠道号
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
func UploadLogToAlyOss(header *multipart.FileHeader, c *gin.Context) (string, error) {
	ossPath := "" //OSS 存储路径

	localfile := "storage/" + header.Filename //本地文件路径
	// gin 简单做了封装,拷贝了文件流
	if err := c.SaveUploadedFile(header, localfile); err != nil {
		logging.Error(err)
		return ossPath, err
	}

	ossPath, ossErr := aliyun.UploadFileToAlyOss(localfile, header.Filename, 1)
	if ossErr != nil {
		logging.Error(ossErr)
		return ossPath, ossErr
	}

	// 删除 log 文件
	os.Remove(localfile)

	return ossPath, nil
}
