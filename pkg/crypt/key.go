package crypt

import (
	"errors"
)

var (
	errorKey       = errors.New("key is nil")
	errorKeyNotSet = errors.New("key do not set to cache")
)

func GetKeyBySdkVersion(sdkVersion string) ([]byte, error) {
	redisKey := "rsa_key:" + sdkVersion
	if x, found := GCache.Get(redisKey); found {
		if key, ok := x.([]byte); ok {
			return key, nil
		}
		return nil, errorKey
	}
	return nil, errorKeyNotSet
}

//func setAllSdkKeyFromRedis() error {
//	dataMap, err := gredis.HGetAll("go:k:prv")
//	if err != nil {
//		return err
//	}
//	if dataMap != nil && *dataMap != nil {
//		for redisField, v := range *dataMap {
//			key := "rsa_key:" + redisField
//			newData := []byte(v)
//			oldData, err := GetKeyBySdkVersion(redisField)
//			if err != nil {
//				GCache.Set(key, newData, cache.NoExpiration)
//				continue
//			}
//			if !reflect.DeepEqual(newData, oldData) {
//				GCache.Set(key, newData, cache.NoExpiration)
//			}
//
//		}
//	}
//
//	return nil
//}
