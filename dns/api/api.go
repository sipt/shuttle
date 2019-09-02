package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/global/namespace"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/dns/cache", listHandleFunc)
	e.DELETE("/api/dns/cache", clearHandleFunc)
}

func listHandleFunc(c *gin.Context) {
	np := namespace.NamespaceWithContext(c)
	cache := np.Profile().DNSCache()
	cacheList := cache.List()
	list := make([]*DNS, 0, len(cacheList))
	for _, v := range cacheList {
		list = append(list, &DNS{
			Typ:         v.Typ,
			Domain:      v.Domain,
			IP:          v.CurrentIP.String(),
			DNSServer:   v.CurrentServer.String(),
			CountryCode: v.CurrentCountry,
		})
	}
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: list,
	})
}
func clearHandleFunc(c *gin.Context) {
	np := namespace.NamespaceWithContext(c)
	cache := np.Profile().DNSCache()
	cache.Clear()
	c.JSON(http.StatusOK, &model.Response{})
}

type DNS struct {
	Typ         string `json:"typ"`
	Domain      string `json:"domain"`
	IP          string `json:"ip"`
	ExpireAt    string `json:"expire_at"`
	DNSServer   string `json:"dns_server"`
	CountryCode string `json:"country_code"`
}
