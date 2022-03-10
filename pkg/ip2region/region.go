package ip2region

import (
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"go-trailer-api/pkg/logging"
)

var IpRegionDb *ip2region.Ip2Region

func Setup() {
	IpRegionDb = createRegionDb()
}

func createRegionDb() *ip2region.Ip2Region {
	regionDb, err := ip2region.New("conf/ip2region.db")

	defer regionDb.Close()
	if err != nil {
		logging.Error(err)
		return nil
	}

	return regionDb
}
