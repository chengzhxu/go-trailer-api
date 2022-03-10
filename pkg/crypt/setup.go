package crypt

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	GCache            *cache.Cache
	refreshKeyLimiter <-chan time.Time
)

func Setup() {
	GCache = cache.New(5*time.Minute, 10*time.Minute)
	refreshCacheKeyData()

	refreshKeyLimiter = time.Tick(10 * time.Minute)
	tickRefreshCacheData()
}

func refreshCacheKeyData() {
	//setAllSdkKeyFromRedis()
}

func tickRefreshCacheData() {
	go func() {
		for {
			<-refreshKeyLimiter
			refreshCacheKeyData()
		}
	}()
}
