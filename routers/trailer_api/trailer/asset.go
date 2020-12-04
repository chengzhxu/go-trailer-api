package trailer

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/gredis"
	"net/http"
)

// @tags Trailer Asset
// @Summary Sync Asset
// @Description 同步 Asset 素材信息
// @ID Sync Asset
// @Produce json
// @Param name body gredis.Asset true "SyncAsset"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/trailer/sync_asset [post]
func SyncTrailerAsset(c *gin.Context) {
	appG := app.Gin{C: c}
	jsonRequest := gredis.Asset{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	if err := jsonRequest.SyncAssetToRedis(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorSyncAssetError, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
