package nacos

import (
	"github.com/goinggo/mapstructure"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go-trailer-api/pkg/setting"
	"gopkg.in/yaml.v2"
	"log"
)

var nacosServer *setting.NacosServer
var nacosConf *setting.NacosConf
var NacosClient config_client.IConfigClient

type TrailerConf struct {
	trailerStatsDB setting.StatsDatabase `yaml:"mysql-stats-db"`
	trailerDB      *setting.TrailerDatabase
	trailerRedis   *setting.Redis
}

func Setup() {
	nacosServer = setting.NacosServerSetting
	nacosConf = setting.NacosConfSetting
	NacosClient = createNacosClient()
	initConf()
}

func createNacosClient() config_client.IConfigClient {
	//服务端配置
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosServer.IpAddr, // nacos服务端的地址, 集群版配置多个
			Port:   nacosServer.Part,   // nacos 的端口
		},
	}

	LogDir := "/tmp/nacos/log"
	CacheDir := "/tmp/nacos/cache"
	if len(nacosServer.LogDir) > 0 {
		LogDir = nacosServer.LogDir
	}
	if len(nacosServer.CacheDir) > 0 {
		CacheDir = nacosServer.CacheDir
	}

	// 客户端配置
	cc := constant.ClientConfig{
		NamespaceId:         nacosServer.NamespaceId, // namespace_id
		TimeoutMs:           10 * 1000,               // http请求超时时间，单位毫秒
		NotLoadCacheAtStart: true,
		LogDir:              LogDir,
		CacheDir:            CacheDir,
		RotateTime:          "1h",
		MaxAge:              3,
		//LogLevel:            "debug",
	}

	nacosClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		log.Fatalf("Stats Nacos.Setup err: %v", err)
	}

	return nacosClient
}

func initConf() {
	confStr, err := NacosClient.GetConfig(vo.ConfigParam{
		DataId: nacosConf.DataId,
		Group:  nacosConf.Group,
	})
	if err != nil {
		log.Fatalf("Stats Nacos.GetConfig err: %v", err)
	}

	confresult := make(map[string]interface{})
	yaml.Unmarshal([]byte(confStr), &confresult)

	statsDbSetting := confresult["mysql-stats-db"].(map[interface{}]interface{})
	mapstructure.Decode(statsDbSetting, &setting.StatsDbSetting)

	trailerDbSetting := confresult["mysql-trailer-db"].(map[interface{}]interface{})
	mapstructure.Decode(trailerDbSetting, &setting.TrailerDbSetting)

	redisDbSetting := confresult["redis-db"].(map[interface{}]interface{})
	mapstructure.Decode(redisDbSetting, &setting.RedisSetting)
}
