package e

const (
	Success       = 1
	Error         = 500
	InvalidParams = 400

	// internal
	ErrorInsertSdkEvent            = 50001
	ErrorInsertDevice              = 50020
	ErrorInsertSdkError            = 50040
	ErrorSyncAssetError            = 60001
	ErrorGetAssetError             = 60010
	ErrorGetAssetEmptyDeviceError  = 60011
	ErrorGetAssetEmptyChannelError = 60012
	ErrorGetAssetEmptyPageError    = 60013
	ErrorGetNewAppError            = 60020

	ErrorEncryptError   = 50050
	ErrorSignatureError = 50051
)
