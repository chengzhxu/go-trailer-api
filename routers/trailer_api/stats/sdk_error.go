package stats

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/service/stats_service"
	"net/http"
)

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
		appG.Response(http.StatusInternalServerError, e.ErrorInsertSdkError, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
