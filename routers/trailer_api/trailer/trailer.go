package trailer_api

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/gredis"
)

// @tags Trailer
// @Summary Get TrailerList
// @Description Get TrailerList
// @ID Get TrailerList
// @Produce json
// @Param TotalRegionUV body gredis.TrailerListParam true "TrailerListParam"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/trailer/get_trailers [post]
func GetTrailers(c *gin.Context) {
	appG := app.Gin{C: c}
	jsonRequest := gredis.TrailerListParam{}
	httpCode, errCode, err := app.BindAndValid(c, &jsonRequest)
	if err != nil {
		appG.Response(httpCode, errCode, err.Error())
		return
	}

	//appG.Response(http.StatusOK, e.Success, data)
}
