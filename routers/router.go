package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go-trailer-api/pkg/middleware"
	"go-trailer-api/pkg/setting"
	"go-trailer-api/routers/bird/user"
	"go-trailer-api/routers/trailer_api/app"
	"go-trailer-api/routers/trailer_api/console"
	"go-trailer-api/routers/trailer_api/stats"
	"go-trailer-api/routers/trailer_api/testing"
	"go-trailer-api/routers/trailer_api/trailer"
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
	apiStats.Use(middleware.CheckToken())
	//apiStats.Use(jwt.JWT())
	{
		//设备信息上报
		apiStats.POST("record_device", stats.InsertDevice)

		//记录 SDK 统计事件
		apiStats.POST("record_sdk_event", stats.InsertSdkEvent)

		//记录 SDK 统计事件 - 加密
		apiStats.POST("record_secret_sdk_event", stats.InsertSecretSdkEvent)

		//SDK 错误信息上报
		apiStats.POST("record_sdk_err", stats.InsertSdkError)

		//APP 应用日志上报
		apiStats.POST("upload_app_log", stats.UploadAppLog)
	}

	apiTrailer := r.Group("/trailer_api/trailer")
	{
		//从后台 同步 Asset 素材信息
		apiTrailer.POST("sync_asset", trailer.SyncTrailerAsset)

		//app 端获取 Asset 素材信息   不加密
		apiTrailer.POST("get_trailer_list", trailer.GetTrailerList)

		//app 端获取 Asset 素材信息   加密
		apiTrailer.POST("get_secret_trailer_list", trailer.GetSecretTrailerList)
	}

	apiApp := r.Group("/trailer_api/app")
	{
		//获取更新 APP 最新版本
		apiApp.POST("get_new_app", app.GetNewAppInfo)

		//获取配置信息
		apiApp.GET("get_trailer_conf", app.GetTrailerConf)
	}

	consoleApp := r.Group("/trailer_api/console")
	{
		//重写 Redis 素材数据 - 根据排序 - display_order
		consoleApp.GET("reset_asset", console.ResetAsset)

		//从 Redis 移除指定素材信息  - 异常数据清理
		consoleApp.GET("remove_asset/:id", console.RemoveAsset)
	}

	testApp := r.Group("/trailer_api/test")
	testApp.Use(middleware.CheckToken())
	{
		//测试 test
		testApp.GET("check_interface", testing.CheckInterface)

		//解密接口 test
		testApp.POST("check_secret_interface", testing.CheckSecretInterface)
	}

	apiBird := r.Group("/bird/userService")
	{
		apiBird.GET("listing", user.Listing) //用户 list
		apiBird.POST("add", user.AddUser)    //新增用户
	}

	return r
}
