package middleware

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/e"
	"net/http"
	"os"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		token := c.Request.Header.Get("token")
		if token == "" {
			appG.Response(http.StatusUnauthorized, e.Unauthorized, nil)
			c.Abort()
		}

		if token != os.Getenv("API_TOKEN") {
			appG.Response(http.StatusUnauthorized, e.AuthorizationError, nil)
			c.Abort()
		}

		c.Next()
	}
}
