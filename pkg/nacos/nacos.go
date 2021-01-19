package nacos

import (
	"fmt"
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
		log.Fatalf("Nacos.Setup err: %v", err)
	}

	return nacosClient
}

func initConf() {
	confStr, err := NacosClient.GetConfig(vo.ConfigParam{
		DataId: nacosConf.DataId,
		Group:  nacosConf.Group,
	})
	if err != nil {
		log.Fatalf("Nacos.GetConfig err: %v", err)
	}

	confresult := make(map[string]interface{})
	yaml.Unmarshal([]byte(confStr), &confresult)

	for k, v := range confresult {
		st := v.(map[interface{}]interface{})
		if _, ok := st["password"]; ok {
			st["password"] = fmt.Sprintf("%v", st["password"])
		}

		mapTo(k, st)
	}
}

func mapTo(section string, v interface{}) {
	switch section {
	case "mysql-stats-db": //MySql trailer_stats
		mapstructure.Decode(v, &setting.StatsDbSetting)
		break
	case "mysql-trailer-db": //MySql trailer
		mapstructure.Decode(v, &setting.TrailerDbSetting)
		break
	case "redis-db": //Redis
		mapstructure.Decode(v, &setting.RedisSetting)
		break
	case "standby-time": //Standby Time
		mapstructure.Decode(v, &setting.StandbyTimeSetting)
		break
	}
}
