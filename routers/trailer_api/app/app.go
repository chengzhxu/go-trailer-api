package app

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/service/app_service"
	"net/http"
)

// @tags APP
// @Summary UPDATE APP
// @Description 获取最新的 APP 版本信息
// @ID UPDATE APP
// @Produce json
// @Param name body app_service.AppParam true "UPDATE_APP"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/app/get_new_app [post]
func GetNewAppInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	jsonRequest := app_service.AppParam{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	appInfo, err := jsonRequest.GetNewAppVersion()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorGetNewAppError, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.Success, appInfo)
}
