package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/ffjson/ffjson"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/service/app_service"
	"go-trailer-api/pkg/util"
	"go-trailer-api/routers/trailer_api"
	"net/http"
)

// @tags APP
// @Summary UPDATE APP
// @Description 获取最新的 APP 版本信息
// @ID UPDATE APP
// @Produce json
// @Param name body model.EData true "UPDATE_APP"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/app/get_new_app [post]
func GetNewAppInfo(c *gin.Context) {
	appG := app.Gin{C: c}

	pData, err := trailer_api.GinDecryptData(c, appG)
	if err != nil {
		logging.Error(err)
		appG.ResponseJson(http.StatusBadRequest, err.Error())
		return
	}

	jsonRequest := app_service.AppParam{}
	jsonRequest.IsSecure = true
	//httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	e := ffjson.Unmarshal(pData.Data, &jsonRequest)
	if e != nil {
		logging.Error(e)
		appG.ResponseJson(http.StatusBadRequest, e.Error())
		return
	}

	appInfo, err := jsonRequest.GetNewAppVersion()
	if err != nil {
		logging.Error(err)
		appG.ResponseJson(http.StatusInternalServerError, err.Error())
		//appG.ResponseEncryptJson(http.StatusInternalServerError,nil, nil)
		return
	}

	res, _ := ffjson.Marshal(appInfo)

	//appG.Response(http.StatusOK, e.Success, appInfo)
	appG.ResponseEncryptJson(http.StatusOK, []byte(res), pData.Key)
}

// @tags Config
// @Summary Trailer Conf
// @Description 获取配置信息
// @ID Trailer Conf
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/app/get_trailer_conf [get]
func GetTrailerConf(c *gin.Context) {
	appG := app.Gin{C: c}
	res := make(map[string]interface{})

	standbyTime := 30
	err, time := util.GetStandbyTime() //待机时长
	if err == nil {
		standbyTime = time
	} else {
		logging.Error(err)
	}

	appPackages := util.GetAppPackage() //app package 下载地址

	res["standby_time"] = standbyTime
	res["app_package"] = appPackages

	appG.Response(http.StatusOK, e.Success, res)
}
