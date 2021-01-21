package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

//Redis
type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = Redis{}

//MySql - Stats
type StatsDatabase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var StatsDbSetting = StatsDatabase{}

//MySql - Trailer
type TrailerDatabase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var TrailerDbSetting = TrailerDatabase{}

//nacos server
type NacosServer struct {
	IpAddr      string
	Part        uint64
	NamespaceId string
	LogDir      string
	CacheDir    string
}

var NacosServerSetting = &NacosServer{}

//nacos config
type NacosConf struct {
	DataId string
	Group  string
}

var NacosConfSetting = &NacosConf{}

//nacos app_package config
type NacosAppPackageConf struct {
	DataId string
	Group  string
}

var NacosAppPackageConfSetting = &NacosAppPackageConf{}

//nacos client config
type NacosClientConf struct {
	DataId string
	Group  string
}

var NacosClientConfSetting = &NacosClientConf{}

//ALiYun OSS
type ALiYunOss struct {
	AccessKeyId           string `mapstructure:"access_key_id"`
	AccessKeySecret       string `mapstructure:"access_key_secret"`
	Endpoint              string `mapstructure:"endpoint"`
	BucketName            string `mapstructure:"bucket_name"`
	ShaFaLauncherPath     string `mapstructure:"shafa_launcher_path"`
	BuDingScreensaverPath string `mapstructure:"buding_screensaver_path"`
}

var ALiYunOssSetting = ALiYunOss{}

//APP 日志白名单
type AppLogWhiteList struct {
	DeviceNo    string `mapstructure:"device_no"`
	ChannelCode string `mapstructure:"channel_code"`
}

var AppLogWhiteListSetting = AppLogWhiteList{}

//客户端待机时长配置
type StandbyTimeConf struct {
	Duration int
}

var StandbyTimeSetting = StandbyTimeConf{}

var cfg *ini.File
var conf string

func Setup() {
	var err error

	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'nacos config': %v", err)
	}

	mapTo("server", ServerSetting)
	//mapTo("redis-db", RedisSetting)
	//mapTo("mysql-stats-db", StatsDbSetting)
	//mapTo("mysql-trailer-db", TrailerDbSetting)
	mapTo("nacos-server", NacosServerSetting)
	mapTo("nacos-config", NacosConfSetting)
	mapTo("nacos-app-package-config", NacosAppPackageConfSetting)
	mapTo("nacos-client-config", NacosClientConfSetting)

	// server
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	// redis
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
