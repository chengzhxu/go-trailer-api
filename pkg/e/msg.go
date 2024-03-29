package e

var MsgFlags = map[int]string{
	Success:            "ok",
	Error:              "fail",
	InvalidParams:      "请求参数错误",
	Unauthorized:       "Token 传入失败",
	AuthorizationError: "无效 Token",

	//
	ErrorInsertSdkEvent:            "记录 SDK 统计事件时失败",
	ErrorValidateSdkEvent:          "SDK 事件统计参数校验失败",
	ErrorCheckSdkEvent:             "SDK 事件统计参数解析失败",
	ErrorInsertDevice:              "记录设备信息时失败",
	ErrorInsertSdkError:            "记录 SDK 错误信息时失败",
	ErrorSyncAssetError:            "同步 Asset 素材信息失败",
	ErrorResetAssetError:           "重置 Asset 素材信息失败",
	ErrorGetAssetError:             "获取 Asset 素材信息失败",
	ErrorGetAssetEmptyDeviceError:  "device_no 不能为空",
	ErrorGetAssetEmptyChannelError: "channel_code 不能为空",
	ErrorGetAssetEmptyPageError:    "page 不能为空",
	ErrorGetNewAppError:            "获取 APP 更新版本失败",
	ErrorGetUploadAppLogError:      "文件接收失败，请检查！",
	ErrorUploadAppLogError:         "文件上传失败！",
	ErrorUploadAppLogToAlyError:    "文件传入阿里云失败！",
	ErrorUploadAppLogTooLargeError: "文件过大，不能超过25M",

	ErrorEncryptError:   "加密失败",
	ErrorSignatureError: "签名验证失败",

	// Bird Users
	ErrorAddUserError: "用户添加失败！",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}
