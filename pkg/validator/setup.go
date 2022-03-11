package validator

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Setup() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic(errors.New("Binding Validator Engine Error\n"))
	}

	arr := map[string]func(fl validator.FieldLevel) bool{
		"bas_date":               BasDate,
		"bas_primary":            BasPrimary,
		"bas_time":               BasTime,
		"nm_bas_time":            NoMustBasTime,
		"int_status":             IntStatus,
		"sdk_event_type":         EventType,
		"sdk_event_kt":           EventKt,
		"sdk_event_info":         EventInfo,
		"obj_sdk_events":         ObjSdkEvents,
		"asset_view_limit":       AssetViewLimit,
		"asset_type":             AssetType,
		"asset_act_type":         AssetActType,
		"asset_score":            AssetScore,
		"asset_is_del":           AssetIsDel,
		"asset_open_apps":        AssetOpenApps,
		"asset_channel_code":     AssetChannelCode,
		"asset_ban_channel_code": AssetBanChannelCode,
		"asset_region_code":      AssetRegionCode,
		"app_log_type":           AppLogType,
		"gender":                 Gender,
	}

	for k, v := range arr {
		if err := validate.RegisterValidation(k, v); err != nil {
			panic(err)
		}
	}
}
