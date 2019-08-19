package dns

import (
	"net"

	"github.com/sirupsen/logrus"

	"github.com/oschwald/geoip2-golang"
	"github.com/sipt/shuttle/assets"
)

var geoipDB *geoip2.Reader

func InitGeoIP(dbFile string) error {
	var err error
	geoipFileBytes, err := assets.ReadFile(dbFile)
	if err != nil {
		return err
	}
	geoipDB, err = geoip2.FromBytes(geoipFileBytes)
	if err != nil {
		return err
	}
	return nil
}

func GeoLookUp(ip net.IP) (countryCode string) {
	if geoipDB == nil {
		return
	}
	country, err := geoipDB.Country(ip)
	if err == nil && country != nil {
		logrus.Debugf("[GeoIP] lookup [%s] country -> [%s]", ip, country.Country.IsoCode)
		return country.Country.IsoCode
	}
	logrus.Errorf("[GeoIP] lookup [%s] country failed: %s", ip, err.Error())
	return
}

func CloseGeoDB() error {
	return geoipDB.Close()
}
