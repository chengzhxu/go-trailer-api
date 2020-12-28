package testing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/model"
	"go-trailer-api/routers/trailer_api"
	"net/http"
)

// @tags Test
// @Summary Test Interface
// @Description 测试接口
// @ID Test
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/test/check_interface [get]
func CheckInterface(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.Success, nil)
}

// @tags Test
// @Summary Test CheckSecretInterface
// @Description APP 更新加密测试接口
// @ID CheckSecretInterface
// @Produce json
// @Param name body model.EDataResponse true "UPDATE_APP"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /trailer_api/test/check_secret_interface [post]
func CheckSecretInterface(c *gin.Context) {
	appG := app.Gin{C: c}

	enReq := &model.EDataResponse{}
	err := c.ShouldBindJSON(enReq)
	if err != nil {
		logging.Info(err)
		return
	}

	res, err := trailer_api.UnClientPack(enReq)
	if err != nil {
		logging.Error(err)
		return
	}

	appG.Response(http.StatusOK, e.Success, fmt.Sprintf("%s", res))
}
