package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/proxy"
)

type Group struct {
	Name       string    `json:"name"`
	Servers    []*Server `json:"servers"`
	SelectType string    `json:"select_type"`
}
type Server struct {
	Name     string `json:"name"`
	Selected bool   `json:"selected"`
	Rtt      string `json:"rtt,omitempty"`
}

func ServerList(ctx *gin.Context) {
	gs := proxy.GetGroups()
	groups := make([]*Group, len(gs))
	var name string
	var group *Group
	for i, g := range gs {
		group = &Group{
			Name:       g.Name,
			Servers:    make([]*Server, len(g.Servers)),
			SelectType: g.SelectType,
		}
		for j, s := range g.Servers {
			name = s.(proxy.IServer).GetName()
			group.Servers[j] = &Server{
				Name:     name,
				Selected: g.Selector.Current().GetName() == name,
			}
			if g.SelectType == "rtt" {
				if ser, ok := s.(*proxy.Server); ok {
					if ser.Rtt == 0 {
						group.Servers[j].Rtt = "failed"
					} else {
						group.Servers[j].Rtt = fmt.Sprintf("%dms", ser.Rtt.Nanoseconds()/1000000)
					}
				}
			}
		}
		groups[i] = group
	}
	ctx.JSON(200, Response{
		Data: groups,
	})
}
func SelectServer(ctx *gin.Context) {
	groupName := ctx.PostForm("group")
	serverName := ctx.PostForm("server")
	if len(groupName) == 0 || len(serverName) == 0 {
		ctx.JSON(500, Response{
			Code: 1, Message: "group or server is empty",
		})
		return
	}
	err := proxy.SelectServer(groupName, serverName)
	if err != nil {
		ctx.JSON(500, Response{
			Code: 1, Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, Response{})
}
func SelectRefresh(ctx *gin.Context) {
	groupName := ctx.PostForm("group")
	if len(groupName) == 0 {
		ctx.JSON(500, Response{
			Code: 1, Message: "group is empty",
		})
		return
	}
	err := proxy.SelectRefresh(groupName)
	if err != nil {
		ctx.JSON(500, Response{
			Code: 1, Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, Response{})
}
