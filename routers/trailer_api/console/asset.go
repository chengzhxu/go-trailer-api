package console

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/gredis"
	"net/http"
	"strconv"
)

// @tags Asset Console
// @Summary Reset Redis Asset
// @Description 清洗 Redis 素材数据
// @ID ResetAsset
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/console/reset_asset [get]
func ResetAsset(c *gin.Context) {
	appG := app.Gin{C: c}

	err := gredis.ResetAsset()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorResetAssetError, err.Error())
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// @tags Asset Console
// @Summary Remove Redis Asset
// @Description 清除 Redis 素材数据
// @ID RemoveAsset
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/console/remove_asset/{id} [get]
func RemoveAsset(c *gin.Context) {
	appG := app.Gin{C: c}
	id, er := strconv.Atoi(c.Param("id"))
	if er != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorResetAssetError, er.Error())
	}

	err := gredis.RemoveAsset(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorResetAssetError, err.Error())
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
