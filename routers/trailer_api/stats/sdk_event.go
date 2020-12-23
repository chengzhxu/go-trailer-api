package stats

import (
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
// @Param name body stats_service.SdkEvent true "Events"
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

	if !util.CheckParamSignature(c) { //检查 签名
		appG.Response(http.StatusInternalServerError, e.ErrorSignatureError, err)
		return
	}

	if jsonRequest.IP == "" {
		r := c.Request
		ip := util.ClientPublicIP(r)
		if ip == "" {
			ip = util.ClientIP(r)
		}
		jsonRequest.IP = ip
	}

	if err := jsonRequest.Insert(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
