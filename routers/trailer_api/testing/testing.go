package testing

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
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
