package shuttle

import (
	"github.com/sipt/shuttle/pool"
	"io"
)

type DirectChannel struct{}

func (d *DirectChannel) Transport(lc, sc IConn) {
	go d.send(sc, lc)
	d.send(lc, sc)
	lc.Close()
	sc.Close()
}

func (d *DirectChannel) send(from, to IConn) {
	var (
		buf = pool.GetBuf()
		n   int
		err error
	)
	for {
		n, err = from.Read(buf)
		if n == 0 {
			return
		}
		if err != nil {
			if err != io.EOF {
				Logger.Errorf("ConnectID [%d] DirectChannel Transport: %v", from.GetID(), err)
			}
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			if err != io.EOF {
				Logger.Error("ConnectID [%d] DirectChannel Transport: %v", to.GetID(), err)
			}
			return
		}
	}
}
