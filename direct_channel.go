package shuttle

import (
	"github.com/sipt/shuttle/pool"
)

type DirectChannel struct{}

func (d *DirectChannel) Transport(lc, sc IConn) {
	go d.send(lc, sc)
	d.send(sc, lc)
}

func (d *DirectChannel) send(from, to IConn) {
	var (
		buf = pool.GetBuf()
		n   int
		err error
	)
	for {
		n, err = from.Read(buf)
		if err != nil {
			from.Close()
			to.Close()
			Logger.Errorf("ConnectID [%d] DirectChannel Transport: %v", from.GetID(), err)
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			from.Close()
			to.Close()
			Logger.Error("ConnectID [%d] DirectChannel Transport: %v", to.GetID(), err)
			return
		}
	}
}
