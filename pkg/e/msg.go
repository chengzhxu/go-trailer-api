package e

var MsgFlags = map[int]string{
	Success:       "ok",
	Error:         "fail",
	InvalidParams: "请求参数错误",

	//
	ErrorInsertSdkEvent: "记录 SDK 统计事件时失败",
	ErrorInsertDevice:   "记录设备信息时失败",
	ErrorInsertSdkError: "记录 SDK 错误信息时失败",
	ErrorSyncAssetError: "同步 Asset 素材信息失败",
	ErrorGetAssetError:  "获取 Asset 素材信息失败",
	ErrorGetNewAppError: "获取 APP 更新版本失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}
