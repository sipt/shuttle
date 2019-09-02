package typ

import "github.com/sipt/shuttle/conn"

type HandleFunc func(conn.ICtxConn)

type K struct {
	V string
}
