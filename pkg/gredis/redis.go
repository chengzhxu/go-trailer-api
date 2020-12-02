package gredis

import (
	"github.com/haxqer/gintools/logging"
	"go-trailer-api/pkg/setting"
	"github.com/gomodule/redigo/redis"
	"time"
)


var RedisConn *redis.Pool

func Setup() {
	RedisConn = createRedisConn(setting.RedisSetting)
}

func createRedisConn(redisSetting *setting.Redis) *redis.Pool {
	redisConn := &redis.Pool{
		MaxIdle:     redisSetting.MaxIdle,
		MaxActive:   redisSetting.MaxActive,
		IdleTimeout: redisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisSetting.Host)
			if err != nil {
				return nil, err
			}
			if redisSetting.Password != "" {
				if _, err := c.Do("AUTH", redisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},

	}

	return redisConn
}

func HGetAll(key string) (*map[string]string, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}


	return &reply, nil
}


func SetString(key string, val string) (bool, error) {
	conn := RedisConn.Get()
	_, err := conn.Do("Set", key, val)
	if err != nil {
		logging.Info(err)
		return false, err
	}

	return true, nil
}

func GetString(key string) (string, error) {
	conn := RedisConn.Get()
	val, err := redis.String(conn.Do("Get", key))
	if err != nil {
		logging.Info(err)
		return "", err
	}

	return val, err
}
