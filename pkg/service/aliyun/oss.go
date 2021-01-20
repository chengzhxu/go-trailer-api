package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go-trailer-api/pkg/setting"
	"go-trailer-api/pkg/util"
	"strings"
)

// 上传文件到阿里云 OSS
func UploadFileToAlyOss(localfile string, uploadfile string, logType int) (string, error) {
	ossPath := "" //存储到阿里云 oss 的地址

	endpoint := setting.ALiYunOssSetting.Endpoint               // Endpoint
	accessKeyId := setting.ALiYunOssSetting.AccessKeyId         // AccessKey
	accessKeySecret := setting.ALiYunOssSetting.AccessKeySecret // Secret
	bucketName := setting.ALiYunOssSetting.BucketName
	// 上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	savePath := ""    //OSS 存储路径
	if logType == 1 { //布丁屏保SDK
		savePath = setting.ALiYunOssSetting.BuDingScreensaverPath
	} else { //沙发桌面
		savePath = setting.ALiYunOssSetting.ShaFaLauncherPath
	}
	if savePath != "" {
		length := len(savePath)
		if util.SubStr(savePath, 0, 1) == "/" {
			savePath = util.SubStr(savePath, 1, length-1)
		}
	}
	objectName := savePath + uploadfile
	// 由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := localfile
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return ossPath, err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return ossPath, err
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return ossPath, err
	}

	ossPath = getOssSavePath(objectName)

	return ossPath, nil
}

// 拼接 OSS 返回路径
func getOssSavePath(savePath string) string {
	endpoint := setting.ALiYunOssSetting.Endpoint // Endpoint
	bucketName := setting.ALiYunOssSetting.BucketName

	bucketNamePath := strings.ReplaceAll(bucketName, "/", "")
	endpointPath := strings.Replace(endpoint, "http://", "", -1)
	endpointPath = strings.Replace(endpointPath, "https://", "", -1)

	return bucketNamePath + "." + endpointPath + "/" + savePath
}
