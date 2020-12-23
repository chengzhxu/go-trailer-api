package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/ffjson/ffjson"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/service/app_service"
	"go-trailer-api/routers/trailer_api"
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

	pData, err := trailer_api.GinDecryptData(c, appG)
	if err != nil {
		logging.Error(err)
		return
	}

	jsonRequest := app_service.AppParam{}
	//httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	e := ffjson.Unmarshal(pData.Data, &jsonRequest)
	if e != nil {
		logging.Error(e)
		appG.ResponseJson(http.StatusBadRequest, nil)
		return
	}

	//appInfo, err := jsonRequest.GetNewAppVersion()
	//if err != nil {
	//	logging.Error(err)
	//	appG.ResponseJson(http.StatusInternalServerError, nil)
	//	//appG.ResponseEncryptJson(http.StatusInternalServerError,nil, nil)
	//	return
	//}

	//res, _ := ffjson.Marshal(appInfo)

	//appG.Response(http.StatusOK, e.Success, appInfo)
	appG.ResponseEncryptJson(http.StatusOK, []byte(`{abcdefggg}`), pData.Key)
}
