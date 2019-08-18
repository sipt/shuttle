package server

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/sipt/shuttle/conn"
)

func init() {
	Register(Direct, func(_, _, _ string, _ map[string]string) (server IServer, e error) {
		return &DirectServer{
			RWMutex: &sync.RWMutex{},
			rtt:     make(map[string]time.Duration),
		}, nil
	})
	Register(Reject, func(_, _, _ string, _ map[string]string) (server IServer, e error) {
		return &RejectServer{}, nil
	})
}

type DirectServer struct {
	rtt map[string]time.Duration
	*sync.RWMutex
}

func (d *DirectServer) Typ() string {
	return Direct
}
func (d *DirectServer) Name() string {
	return Direct
}
func (d *DirectServer) SetRtt(key string, rtt time.Duration) {
	d.Lock()
	defer d.Unlock()
	d.rtt[key] = rtt
}
func (d *DirectServer) Rtt(key string) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.rtt[key]
}
func (d *DirectServer) DialTCP(ctx context.Context, addr, port string, dial conn.DialTCPFunc) (*net.TCPConn, error) {
	return dial(ctx, addr, port)
}

func (d *DirectServer) DialUDP(ctx context.Context, addr, port string, dial conn.DialUDPFunc) (*net.UDPConn, error) {
	return dial(ctx, addr, port)
}

type RejectServer struct{}

func (r *RejectServer) Typ() string {
	return Reject
}
func (r *RejectServer) Name() string {
	return Reject
}
func (r *RejectServer) SetRtt(_ string, _ time.Duration) {
}
func (r *RejectServer) Rtt(_ string) time.Duration {
	return time.Duration(-1)
}
func (r *RejectServer) DialTCP(ctx context.Context, addr, port string, dial conn.DialTCPFunc) (*net.TCPConn, error) {
	return nil, ErrRejected
}

func (r *RejectServer) DialUDP(ctx context.Context, addr, port string, dial conn.DialUDPFunc) (*net.UDPConn, error) {
	return nil, ErrRejected
}
