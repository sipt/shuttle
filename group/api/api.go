package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/events"
	"github.com/sipt/shuttle/global/namespace"
	"github.com/sipt/shuttle/group"
	"github.com/sirupsen/logrus"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/groups", listHandleFunc)
	e.GET("/api/group", groupHandleFunc)
	e.PUT("/api/group/rtt", resetHandleFunc)
	e.PUT("/api/group/select", selectHandleFunc)

	e.GET("/ws/rtt", rttEventsHandleFunc)
}

func listHandleFunc(c *gin.Context) {
	np := namespace.NamespaceWithContext(c)
	groups := np.Profile().Group()
	list := make([]*Group, 0, len(groups))
	if np.Mode() != constant.ModeDirect {
		for _, v := range groups {
			if v.Name() == group.Global && np.Mode() != constant.ModeGlobal {
				continue
			}
			list = append(list, makeGroupResp(v))
		}
		sort.Sort(SortableGroups(list))
	}
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: list,
	})
}

func groupHandleFunc(c *gin.Context) {
	np := namespace.NamespaceWithContext(c)
	groups := np.Profile().Group()
	name := c.Query("group")
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

func selectHandleFunc(c *gin.Context) {
	np := namespace.NamespaceWithContext(c)
	groups := np.Profile().Group()
	name := c.Query("group")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, &model.Response{
			Code:    1,
			Message: "group name is empty",
		})
		return
	}
	subName := c.Query("server")
	if len(subName) == 0 {
		c.JSON(http.StatusBadRequest, &model.Response{
			Code:    1,
			Message: "group sub name is empty",
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
	err := g.Select(subName)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("[%s] not found in group[%s]", subName, name),
		})
	}
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: makeGroupResp(g),
	})
}

func resetHandleFunc(c *gin.Context) {
	np := namespace.NamespaceWithContext(c)
	groups := np.Profile().Group()
	name := c.Query("group")
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
			Name: g.Selected().Name(),
			Typ:  g.Selected().Typ(),
			RTT:  formatRtt(g.Server().Rtt(g.Name())),
		},
		Servers: make([]Server, len(g.Items())),
	}

	for i, s := range g.Items() {
		resp.Servers[i] = Server{
			Name:     s.Name(),
			Typ:      s.Typ(),
			RTT:      formatRtt(s.Server().Rtt(resp.Name)),
			Selected: g.Selected().Name() == s.Name(),
		}
	}
	return resp
}

type SortableGroups []*Group

func (sg SortableGroups) Len() int {
	return len(sg)
}
func (sg SortableGroups) Less(i, j int) bool {
	return sg[i].Name < sg[j].Name
}
func (sg SortableGroups) Swap(i, j int) {
	sg[i], sg[j] = sg[j], sg[i]
}

type Group struct {
	Name     string   `json:"name"`
	Typ      string   `json:"typ"`
	Selected Server   `json:"selected"`
	Servers  []Server `json:"servers"`
}

type Server struct {
	Name     string `json:"name"`
	Typ      string `json:"typ"`
	RTT      string `json:"rtt"`
	Selected bool   `json:"selected"`
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

var (
	wsConnID  int64 = 0
	wsConnMap       = make(map[int64]*websocket.Conn)
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:   2048,
	WriteBufferSize:  2048,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func rttEventsHandleFunc(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	id := atomic.AddInt64(&wsConnID, 1)
	go func() {
		defer conn.Close()
		for {
			typ, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if typ == websocket.CloseMessage {
				delete(wsConnMap, id)
				return
			}
		}
	}()
	wsConnMap[id] = conn
}

func init() {
	events.RegisterEvent(events.GroupRttEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(string)
		if !ok {
			return errors.Errorf("[%s] is not GroupName", reflect.TypeOf(v).Kind().String())
		}
		notifyClient(ctx, r)
		return nil
	})
}

func notifyClient(ctx context.Context, name string) {
	if len(wsConnMap) == 0 {
		return
	}
	np := namespace.NamespaceWithContext(ctx)
	groups := np.Profile().Group()
	g, ok := groups[name]
	if !ok || g == nil {
		return
	}
	resp := makeGroupResp(g)
	var err error
	for _, conn := range wsConnMap {
		err = conn.WriteJSON(resp)
		if err != nil {
			logrus.WithError(err).Error("[group] [rtt.notifyClient] failed")
		}
	}
}
