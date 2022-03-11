package user

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/service/bird/userService"
	"net/http"
)

// @tags Bird
// @Summary User List
// @Description 用户列表
// @ID UserListing
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /bird/userService/listing [get]
func Listing(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.Success, nil)
}

// @tags Bird
// @Summary Add User
// @Description 新增用户
// @ID AddUser
// @Produce json
// @Param name body userService.User true "UPDATE_APP"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /bird/userService/add [post]
func AddUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	jsonRequest := userService.User{}

	if err := c.ShouldBind(&jsonRequest); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, err.Error())
		//c.JSON(400, gin.H{ "error": err.Error() })
		return
	}

	//保存 MySql
	if err := jsonRequest.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAddUserError, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
