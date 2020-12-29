package console

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/gredis"
	"net/http"
)

// @tags Reset Redis
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
