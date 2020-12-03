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

	apiStats := r.Group("/trailer_api/stats")
	//apiStats.Use(jwt.JWT())
	{
		//设备信息上报
		apiStats.POST("record_device", stats.InsertDevice)

		//记录 SDK 统计事件
		apiStats.POST("record_sdk_event", stats.InsertSdkEvent)

		//SDK 错误信息上报
		apiStats.POST("record_sdk_err", stats.InsertSdkError)
	}

	return r
}
