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

var RedisSetting = &Redis{}

//MySql - Stats
type StatsDatabase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var StatsDbSetting = &StatsDatabase{}

//MySql - Trailer
type TrailerDatabase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var TrailerDbSetting = &TrailerDatabase{}

//ALiYun OSS
type ALiYunOss struct {
	AccessKeyId           string
	AccessKeySecret       string
	Endpoint              string
	BucketName            string
	ShaFaLauncherPath     string
	BuDingScreensaverPath string
}

var ALiYunOssSetting = &ALiYunOss{}

//APP 日志白名单
type AppLogWhiteList struct {
	DeviceNo    string
	ChannelCode string
}

var AppLogWhiteListSetting = &AppLogWhiteList{}

//客户端待机时长配置
type StandbyTimeConf struct {
	Duration int
}

var StandbyTimeSetting = &StandbyTimeConf{}

var cfg *ini.File

func Setup() {
	var err error

	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("server", ServerSetting)
	mapTo("redis-db", RedisSetting)
	mapTo("mysql-stats-db", StatsDbSetting)
	mapTo("mysql-trailer-db", TrailerDbSetting)
	mapTo("aliyun-oss", ALiYunOssSetting)
	mapTo("app-log-white-list", AppLogWhiteListSetting)
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
