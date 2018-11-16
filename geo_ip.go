package shuttle

import (
	"github.com/oschwald/geoip2-golang"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/util"
	"net"
)

var geoipDB *geoip2.Reader

type GeoIP struct {
	IP      *net.IPNet
	Country string
}

func InitGeoIP(geoipDBFile string) error {
	var err error
	geoipDB, err = geoip2.Open(geoipDBFile)
	if err != nil {
		return err
	}
	return nil
}

func GeoLookUp(ip net.IP) string {
	country, err := geoipDB.Country(ip)
	if err == nil && country != nil {
		log.Logger.Debugf("[GeoIP] lookup [%s] country -> [%s]", ip.String(), country.Country.IsoCode)
		return country.Country.IsoCode
	}
	r, err := util.WatchIP(ip.String())
	if err == nil && r != nil {
		log.Logger.Debugf("[GeoIP] use taobao api [%s] country -> [%s]", ip.String(), r.CountryID)
		return r.CountryID
	}
	return ""
}

func CloseGeoDB() error {
	return geoipDB.Close()
}
