package gredis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/setting"
	"strconv"
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

func Zadd(key string) bool {
	conn := RedisConn.Get()

	for i := 1; i < 50; i++ {
		_, err := conn.Do("zadd", key, i, "asset"+strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
		}
	}

	return true
}

func Zrem(key string) bool {
	conn := RedisConn.Get()

	_, err := conn.Do("zrem", key, "asset10")
	if err != nil {
		fmt.Println(err)
	}

	return true
}

func ZrByScore(key string) bool {
	conn := RedisConn.Get()

	res, err := redis.Values(conn.Do("zrevrangebyscore", key, 7889155200, 1))
	if err != nil {
		fmt.Println(err)
	}

	assetArray := AssetArray{}
	for _, v := range res {
		reply, err := redis.Bytes(conn.Do("hget", "trailer_asset", v))
		if err != nil {
			fmt.Printf("%s ", v.([]byte))
		}
		var asset *Asset

		json.Unmarshal(reply, &asset)
		fmt.Println(asset.Id)
		assetArray = append(assetArray, asset)
	}

	fmt.Printf("%s ", assetArray)

	return true
}

func ZCard(key string) bool {
	conn := RedisConn.Get()

	num, err := conn.Do("zcard", key)
	if err != nil {
		fmt.Println(err)
	}
	//
	//for _, v := range res {
	//	fmt.Printf("%s ", v)
	//}

	fmt.Println(num)

	return true
}
