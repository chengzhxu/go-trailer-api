package trailer

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/ffjson/ffjson"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/gredis"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/util"
	"go-trailer-api/routers/trailer_api"
	"net/http"
)

// @tags Trailer
// @Summary Get TrailerList
// @Description 获取预告片信息 - 不加密
// @ID Get TrailerList
// @Produce json
// @Param TrailerListParam body gredis.TrailerListParam true "TrailerListParam"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/trailer/get_trailer_list [post]
func GetTrailerList(c *gin.Context) {
	appG := app.Gin{C: c}
	jsonRequest := gredis.TrailerListParam{}
	jsonRequest.IsSecure = true
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	if jsonRequest.RegionCode == "" {
		ipInfo := util.GetIpInfoByRequest(c.Request) //根据请求获取地域信息
		jsonRequest.RegionCode = model.ReadRegionFromFile(ipInfo.City)
	}

	assetRes, err := jsonRequest.QueryTrailerList()
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ErrorGetAssetError, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, assetRes)
}

// @tags Trailer
// @Summary Get SecretTrailerList
// @Description 获取预告片信息 - 加密
// @ID Get SecretTrailerList
// @Produce json
// @Param TrailerListParam body model.EData true "EDataParam"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/trailer/get_secret_trailer_list [post]
func GetSecretTrailerList(c *gin.Context) {
	appG := app.Gin{C: c}

	pData, err := trailer_api.GinDecryptData(c, appG)
	if err != nil {
		logging.Error(err)
		return
	}

	jsonRequest := gredis.TrailerListParam{}
	jsonRequest.IsSecure = true
	//httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	e := ffjson.Unmarshal(pData.Data, &jsonRequest)
	if e != nil {
		logging.Error(e)
		appG.ResponseJson(http.StatusBadRequest, nil)
		return
	}

	ipInfo := util.GetIpInfoByRequest(c.Request) //根据请求获取地域信息
	jsonRequest.RegionCode = model.ReadRegionFromFile(ipInfo.City)

	assetRes, err := jsonRequest.QueryTrailerList()
	if err != nil {
		logging.Error(err)
		appG.ResponseJson(http.StatusInternalServerError, nil)
		//appG.ResponseEncryptJson(http.StatusInternalServerError,nil, nil)
		return
	}

	res, _ := ffjson.Marshal(assetRes)

	//appG.Response(http.StatusOK, e.Success, appInfo)
	appG.ResponseEncryptJson(http.StatusOK, []byte(res), pData.Key)

}
