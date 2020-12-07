package trailer

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/gredis"
	"go-trailer-api/pkg/logging"
	"net/http"
)

// @tags Trailer
// @Summary Get TrailerList
// @Description 获取预告片信息
// @ID Get TrailerList
// @Produce json
// @Param TrailerListParam body gredis.TrailerListParam true "TrailerListParam"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/trailer/get_trailer_list [post]
func GetTrailerList(c *gin.Context) {
	appG := app.Gin{C: c}
	jsonRequest := gredis.TrailerListParam{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	assetRes, err := jsonRequest.QueryTrailerList()
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ErrorGetAssetError, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, assetRes)
}
