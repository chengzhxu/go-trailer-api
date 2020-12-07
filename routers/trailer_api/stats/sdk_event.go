package stats

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/service/stats_service"
	"go-trailer-api/pkg/util"
	"net/http"
)

// @tags Stats
// @Summary SDK 事件统计
// @Description SDK 事件统计   参数：事件 json 数组
// @ID Insert SdkEvent
// @Produce json
// @Param name body stats_service.ObjSdkEvents true "Events"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/stats/record_sdk_event [post]
func InsertSdkEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		err  error
	)

	jsonRequest := stats_service.ObjSdkEvents{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	var sdkEvents []stats_service.SdkEvent
	json.Unmarshal([]byte(jsonRequest.SdkEvents), &sdkEvents)

	r := c.Request
	ip := util.ClientPublicIP(r)
	if ip == "" {
		ip = util.ClientIP(r)
	}
	//validator := govalidators.New()

	for _, se := range sdkEvents {
		se.IP = ip
		//errList := validator.Validate(se)
		//if errList != nil {
		//	for _, el := range errList {
		//		appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, nil)
		//		return
		//	}
		//}
		if err := se.Insert(); err != nil {
			appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, err)
			return
		}
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
