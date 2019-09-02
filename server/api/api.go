package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/global"
	"github.com/sipt/shuttle/server"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/servers", func(c *gin.Context) {
		np := global.NamespaceWithContext(c)
		servers := np.Profile().Server()
		list := make([]*ItemResponse, 0, len(servers))
		for _, v := range servers {
			list = append(list, &ItemResponse{
				Name: v.Name(),
				Typ:  v.Typ(),
				RTT:  formatRtt(v.Rtt(server.DefaultRttKey)),
			})
		}
		c.JSON(http.StatusOK, &model.Response{
			Code: 0,
			Data: list,
		})
	})
	e.GET("/api/servers/:name", func(c *gin.Context) {
		np := global.NamespaceWithContext(c)
		servers := np.Profile().Server()
		name := c.Param("name")
		if len(name) == 0 {
			c.JSON(http.StatusBadRequest, &model.Response{
				Code:    1,
				Message: "server name is empty",
			})
			return
		}
		s, ok := servers[name]
		if !ok || s == nil {
			c.JSON(http.StatusBadRequest, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("server name[%s] not found", name),
			})
			return
		}

		c.JSON(http.StatusOK, &model.Response{
			Data: &ItemResponse{
				Name: s.Name(),
				Typ:  s.Typ(),
				RTT:  formatRtt(s.Rtt(server.DefaultRttKey)),
			},
		})
	})
	e.PUT("/api/servers/:name/rtt", func(c *gin.Context) {
		np := global.NamespaceWithContext(c)
		servers := np.Profile().Server()
		name := c.Param("name")
		if len(name) == 0 {
			c.JSON(http.StatusBadRequest, &model.Response{
				Code:    1,
				Message: "server name is empty",
			})
			return
		}
		s, ok := servers[name]
		if !ok || s == nil {
			c.JSON(http.StatusBadRequest, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("server name[%s] not found", name),
			})
			return
		}

		c.JSON(http.StatusOK, &model.Response{
			Data: &ItemResponse{
				Name: s.Name(),
				Typ:  s.Typ(),
				RTT:  formatRtt(s.TestRtt(server.DefaultRttKey, "")),
			},
		})
	})
}

func formatRtt(t time.Duration) string {
	if t > 0 {
		t = t.Round(time.Millisecond)
		return t.String()
	} else if t == 0 {
		return "no rtt"
	} else {
		return "failed"
	}
}

type ItemResponse struct {
	Name string `json:"name"`
	Typ  string `json:"typ"`
	RTT  string `json:"rtt"`
}
