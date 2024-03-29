package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/ffjson/ffjson"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/service/stats_service"
	"go-trailer-api/pkg/util"
	"go-trailer-api/routers/trailer_api"
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

	param := stats_service.MapSdkEvent(jsonRequest)
	if !util.CheckParamSignature(param, jsonRequest.Signature) { //检查 签名
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

	//根据事件参数拆分为上报数据
	eventArr := stats_service.TranceEventKtJson(jsonRequest)
	for _, event := range eventArr {
		if err := event.Insert(); err != nil {
			appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, err)
			return
		}
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// @tags Stats
// @Summary SDK 事件统计-加密
// @Description SDK 事件统计-加密
// @ID Insert Secret SdkEvent
// @Produce json
// @Param name body model.EData true "Events"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/stats/record_secret_sdk_event [post]
func InsertSecretSdkEvent(c *gin.Context) {
	appG := app.Gin{C: c}

	pData, err := trailer_api.GinDecryptData(c, appG)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ErrorValidateSdkEvent, err)
		return
	}

	jsonRequest := stats_service.SdkEvent{}
	er := ffjson.Unmarshal(pData.Data, &jsonRequest)
	if er != nil {
		logging.Error(er)
		appG.Response(http.StatusInternalServerError, e.ErrorCheckSdkEvent, er)
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

	//根据事件参数拆分为上报数据
	eventArr := stats_service.TranceEventKtJson(jsonRequest)
	for _, event := range eventArr {
		if err := event.Insert(); err != nil {
			appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkEvent, err)
			return
		}
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
