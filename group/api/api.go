package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/global"
	"github.com/sipt/shuttle/group"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/groups", listHandleFunc)
	e.GET("/api/groups/:name", groupHandleFunc)
	e.PUT("/api/groups/:name/rtt", resetHandleFunc)
}

func listHandleFunc(c *gin.Context) {
	np := global.NamespaceWithContext(c)
	groups := np.Profile().Group()
	list := make([]*Group, 0, len(groups))
	for _, v := range groups {
		list = append(list, makeGroupResp(v))
	}
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: list,
	})
}

func groupHandleFunc(c *gin.Context) {
	np := global.NamespaceWithContext(c)
	groups := np.Profile().Group()
	name := c.Param("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, &model.Response{
			Code:    1,
			Message: "group name is empty",
		})
		return
	}
	g, ok := groups[name]
	if !ok || g == nil {
		c.JSON(http.StatusBadRequest, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("group name[%s] not found", name),
		})
		return
	}
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: makeGroupResp(g),
	})
}

func resetHandleFunc(c *gin.Context) {
	np := global.NamespaceWithContext(c)
	groups := np.Profile().Group()
	name := c.Param("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, &model.Response{
			Code:    1,
			Message: "group name is empty",
		})
		return
	}
	g, ok := groups[name]
	if !ok || g == nil {
		c.JSON(http.StatusBadRequest, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("group name[%s] not found", name),
		})
		return
	}
	g.Reset()
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: makeGroupResp(g),
	})
}

func makeGroupResp(g group.IGroup) *Group {
	resp := &Group{
		Name: g.Name(),
		Typ:  g.Typ(),
		Selected: Server{
			Name: g.Server().Name(),
			Typ:  g.Server().Typ(),
			RTT:  formatRtt(g.Server().Rtt(g.Name())),
		},
		Servers: make([]Server, len(g.Items())),
	}

	for i, s := range g.Items() {
		resp.Servers[i] = Server{
			Name: s.Name(),
			Typ:  s.Typ(),
			RTT:  formatRtt(s.Server().Rtt(resp.Name)),
		}
	}
	return resp
}

type Group struct {
	Name     string   `json:"name"`
	Typ      string   `json:"typ"`
	Selected Server   `json:"selected"`
	Servers  []Server `json:"servers"`
}

type Server struct {
	Name string `json:"name"`
	Typ  string `json:"typ"`
	RTT  string `json:"rtt"`
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
