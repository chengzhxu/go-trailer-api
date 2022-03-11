package e

const (
	Success            = 1
	Error              = 500
	InvalidParams      = 400
	Unauthorized       = 401
	AuthorizationError = 403

	// internal
	ErrorInsertSdkEvent            = 50001
	ErrorValidateSdkEvent          = 50002
	ErrorCheckSdkEvent             = 50003
	ErrorInsertDevice              = 50020
	ErrorInsertSdkError            = 50040
	ErrorSyncAssetError            = 60001
	ErrorResetAssetError           = 60002
	ErrorGetAssetError             = 60010
	ErrorGetAssetEmptyDeviceError  = 60011
	ErrorGetAssetEmptyChannelError = 60012
	ErrorGetAssetEmptyPageError    = 60013
	ErrorGetNewAppError            = 60020
	ErrorGetUploadAppLogError      = 70001
	ErrorUploadAppLogError         = 70002
	ErrorUploadAppLogToAlyError    = 70003
	ErrorUploadAppLogTooLargeError = 70004

	ErrorEncryptError   = 50050
	ErrorSignatureError = 50051

	//Bird User
	ErrorAddUserError = 80001
)
