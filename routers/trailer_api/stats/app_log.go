package stats

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/service/stats_service"
	"net/http"
)

// @tags Stats
// @Summary APP 应用日志上报
// @Description APP 应用日志上报
// @ID Upload APP LOG
// @Produce json
// @Param name body stats_service.AppLog true "AppLog"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/stats/upload_app_log [post]
func UploadAppLog(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		err  error
	)

	header, FileErr := c.FormFile("log_file")
	if FileErr != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorGetUploadAppLogError, FileErr)
		return
	}

	jsonRequest := stats_service.AppLog{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)

	if err != nil {
		appG.Response(httpCode, errCode, err)
		return
	}

	logUrl, ossErr := stats_service.UploadLogToAlyOss(header, c)
	if ossErr != nil || logUrl == "" {
		appG.Response(http.StatusInternalServerError, e.ErrorUploadAppLogToAlyError, err)
		return
	}
	jsonRequest.URL = logUrl

	//保存 Log MySql
	if err := jsonRequest.Insert(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorUploadAppLogError, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
