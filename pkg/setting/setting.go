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

//客户端待机时长配置
type StandbyTimeConf struct {
	Duration int
}

var StandbyTimeSetting = &StandbyTimeConf{}

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
	mapTo("nacos-trailer", NacosConfSetting)
	mapTo("standby-time", StandbyTimeSetting)

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
