package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle"
)

type Group struct {
	Name    string    `json:"name"`
	Servers []*Server `json:"servers"`
}
type Server struct {
	Name     string `json:"name"`
	Selected bool   `json:"selected"`
}

func ServerList(ctx *gin.Context) {
	gs := shuttle.GetGroups()
	groups := make([]*Group, len(gs))
	var name string
	var group *Group
	for i, g := range gs {
		group = &Group{
			Name:    g.Name,
			Servers: make([]*Server, len(g.Servers)),
		}
		for j, s := range g.Servers {
			name = s.(shuttle.IServer).GetName()
			group.Servers[j] = &Server{
				Name:     name,
				Selected: g.Selector.Current().GetName() == name,
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
	err := shuttle.SelectServer(groupName, serverName)
	if err != nil {
		ctx.JSON(500, Response{
			Code: 1, Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, Response{})
}
