package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/global/namespace"
	"github.com/sipt/shuttle/server"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/servers", func(c *gin.Context) {
		np := namespace.NamespaceWithContext(c)
		servers := np.Profile().Server()
		list := make([]*ItemResponse, 0, len(servers))
		for _, v := range servers {
			list = append(list, &ItemResponse{
				Name: v.Name(),
				Typ:  v.Typ(),
				RTT:  formatRtt(v.Rtt(server.DefaultRttKey)),
			})
		}
		c.JSON(http.StatusOK, &model.Response[[]*ItemResponse]{
			Code: 0,
			Data: list,
		})
	})
	e.GET("/api/servers/:name", func(c *gin.Context) {
		np := namespace.NamespaceWithContext(c)
		servers := np.Profile().Server()
		name := c.Param("name")
		if len(name) == 0 {
			c.JSON(http.StatusBadRequest, &model.Response[any]{
				Code:    1,
				Message: "server name is empty",
			})
			return
		}
		s, ok := servers[name]
		if !ok || s == nil {
			c.JSON(http.StatusBadRequest, &model.Response[any]{
				Code:    1,
				Message: fmt.Sprintf("server name[%s] not found", name),
			})
			return
		}

		c.JSON(http.StatusOK, &model.Response[*ItemResponse]{
			Data: &ItemResponse{
				Name: s.Name(),
				Typ:  s.Typ(),
				RTT:  formatRtt(s.Rtt(server.DefaultRttKey)),
			},
		})
	})
	e.PUT("/api/servers/:name/rtt", func(c *gin.Context) {
		np := namespace.NamespaceWithContext(c)
		servers := np.Profile().Server()
		name := c.Param("name")
		if len(name) == 0 {
			c.JSON(http.StatusBadRequest, &model.Response[any]{
				Code:    1,
				Message: "server name is empty",
			})
			return
		}
		s, ok := servers[name]
		if !ok || s == nil {
			c.JSON(http.StatusBadRequest, &model.Response[any]{
				Code:    1,
				Message: fmt.Sprintf("server name[%s] not found", name),
			})
			return
		}

		c.JSON(http.StatusOK, &model.Response[*ItemResponse]{
			Data: &ItemResponse{
				Name: s.Name(),
				Typ:  s.Typ(),
				RTT:  formatRtt(s.TestRtt(server.DefaultRttKey, "")),
			},
		})
	})
}

func formatRtt(t time.Duration) int64 {
	if t > 0 {
		return t.Milliseconds()
	} else if t == 0 {
		return 0
	} else {
		return -1
	}
}

type ItemResponse struct {
	Name string `json:"name"`
	Typ  string `json:"typ"`
	RTT  int64  `json:"rtt"`
}
