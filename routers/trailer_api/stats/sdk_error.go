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
		appG.Response(httpCode, errCode, err)
		return
	}

	if jsonRequest.Signature != "" { //检查 签名
		param := stats_service.MapSdkError(jsonRequest)
		if !util.CheckParamSignature(param, jsonRequest.Signature) {
			appG.Response(http.StatusInternalServerError, e.ErrorSignatureError, err)
			return
		}
	}

	if err := jsonRequest.Insert(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkError, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
