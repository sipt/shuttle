package dns

import (
	"github.com/oschwald/geoip2-golang"
	"github.com/sipt/shuttle/log"
	"net"
)

var geoipDB *geoip2.Reader

func InitGeoIP(dbFiled string) error {
	var err error
	geoipDB, err = geoip2.Open(dbFiled)
	if err != nil {
		return err
	}
	return nil
}

func GeoLookUp(ip string) (countryCode string) {
	if geoipDB == nil {
		return
	}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return
	}
	country, err := geoipDB.Country(netIP)
	if err == nil && country != nil {
		log.Logger.Debugf("[GeoIP] lookup [%s] country -> [%s]", ip, country.Country.IsoCode)
		return country.Country.IsoCode
	}
	log.Logger.Debugf("[GeoIP] lookup [%s] country failed: %s", ip, err.Error())
	return
}

func CloseGeoDB() error {
	return geoipDB.Close()
}
