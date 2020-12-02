package app

import (
	"go-trailer-api/pkg/e"
	"github.com/gin-gonic/gin"
	"github.com/haxqer/gintools/logging"
	"net/http"
)

type requestInfo struct {
	Request interface{}
	Err string
}

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, b interface{}) (int, int, error) {
	err := c.ShouldBindJSON(b)
	if err != nil {
		logging.Info(&requestInfo{
			Request: b,
			Err:     err.Error(),
		})
		//MarkErrors(err)
		return http.StatusBadRequest, e.InvalidParams, err
	}

	return http.StatusOK, e.Success, nil
}

