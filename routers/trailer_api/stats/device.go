package stats

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/service/stats_service"
	"net/http"
)

// @tags Stats
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
		appG.Response(httpCode, errCode, err)
		return
	}

	if err := jsonRequest.Insert(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorInsertDevice, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
