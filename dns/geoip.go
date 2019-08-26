package dns

import (
	"net"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/oschwald/geoip2-golang"
	"github.com/sipt/shuttle/assets"
)

var fileName = "/Users/sipt/workspace/go/shuttle/GeoLite2-Country.mmdb"

var geoipDB *geoip2.Reader

func InitGeoIP() error {
	var err error
	geoipFileBytes, err := assets.ReadFile(fileName)
	if err != nil {
		return errors.Errorf("reade geo file [%s] failed: %s", fileName, err.Error())
	}
	geoipDB, err = geoip2.FromBytes(geoipFileBytes)
	if err != nil {
		return errors.Errorf("reade geo file [%s] failed: %s", fileName, err.Error())
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
