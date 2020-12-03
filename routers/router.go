package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go-trailer-api/pkg/setting"
	"go-trailer-api/routers/trailer_api/stats"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if setting.ServerSetting.RunMode == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api_stats := r.Group("/trailer_api/stats")
	//api_stats.Use(jwt.JWT())
	{
		//记录 SDK 统计事件
		api_stats.POST("insert_sdk_event", stats.InsertSdkEvent)
	}

	return r
}
