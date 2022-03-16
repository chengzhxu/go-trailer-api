package user

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/model/bird/userModel"
	"go-trailer-api/pkg/service/bird/userService"
	"net/http"
)

// @tags Bird
// @Summary User Listing
// @Description 用户列表
// @ID UserListing
// @Produce json
// @Param name body userModel.UserListParams true "User List"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /bird/user/listing [post]
func Listing(c *gin.Context) {
	appG := app.Gin{C: c}

	jsonRequest := userModel.UserListParams{}

	if err := c.ShouldBind(&jsonRequest); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, err.Error())
		return
	}

	if jsonRequest.PageSize == 0 {
		jsonRequest.PageSize = 20
	}
	if jsonRequest.Page == 0 {
		jsonRequest.Page = 1
	}

	//获取列表信息
	err, list, total := userService.Listing(jsonRequest)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAddUserError, err)
		return
	}

	appG.Response(http.StatusOK, e.Success, app.ListingResult{
		List:     list,
		Total:    total,
		Page:     jsonRequest.Page,
		PageSize: jsonRequest.PageSize,
	})
}

// @tags Bird
// @Summary Add User
// @Description 新增用户
// @ID AddUser
// @Produce json
// @Param name body userService.User true "Add User"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /bird/user/add [post]
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
