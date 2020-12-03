package stats

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/service/stats_service"
	"go-trailer-api/pkg/tool"
	"net/http"
)

// @tags SdkEvent
// @Summary SDK 事件统计
// @Description SDK 事件统计
// @ID Insert SdkEvent
// @Produce json
// @Param name body stats_service.SdkEvent true "Event"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/stats/record_sdk_event [post]
func InsertSdkEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		err  error
	)

	jsonRequest := stats_service.SdkEvent{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	r := c.Request
	ip := tool.ClientPublicIP(r)
	if ip == "" {
		ip = tool.ClientIP(r)
	}
	jsonRequest.IP = ip
	if err := jsonRequest.Insert(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// @tags Device
// @Summary 设备信息上报
// @Description 设备信息上报
// @ID Record Device
// @Produce json
// @Param name body stats_service.Device true "Device"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/stats/record_device [post]
func InsertDevice(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		err  error
	)

	jsonRequest := stats_service.Device{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	if err := jsonRequest.Insert(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// @tags SdkError
// @Summary SDK 错误信息上报
// @Description SDK 错误信息上报
// @ID Record SdkError
// @Produce json
// @Param name body stats_service.SdkError true "SdkError"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/stats/record_sdk_err [post]
func InsertSdkError(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		err  error
	)

	jsonRequest := stats_service.SdkError{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	if err := jsonRequest.Insert(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
