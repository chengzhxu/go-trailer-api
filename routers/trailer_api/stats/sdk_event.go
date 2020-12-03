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
// @Summary Insert SdkEvent
// @Description Insert SdkEvent
// @ID Insert SdkEvent
// @Produce json
// @Param name body stats_service.SdkEvent true "Event"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/stats/insert_sdk_event [post]
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
