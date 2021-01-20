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

	if err := c.ShouldBind(&jsonRequest); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, err.Error())
		//c.JSON(400, gin.H{ "error": err.Error() })
		return
	}

	// too large  限制文件大小
	maxSize := 1024 * 1024 * 21 // 21M
	if header.Size > int64(maxSize) {
		appG.Response(http.StatusInternalServerError, e.ErrorUploadAppLogTooLargeError, err)
		return
	}

	// 上传文件到 OSS
	logUrl, err := stats_service.UploadLogToAlyOss(header, jsonRequest.LogType, c)
	if err != nil || logUrl == "" {
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
