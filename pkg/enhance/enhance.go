package enhance

import (
	"context"
	"sync"

	"github.com/sipt/shuttle/pkg/dns"
	mockdns "github.com/sipt/shuttle/pkg/dns/mock"
	"github.com/sipt/shuttle/pkg/tun"
)

type EnhanceModeState int

const (
	EnhanceModeStateInit EnhanceModeState = iota
	EnhanceModeStateRunning
	EnhanceModeStateStopped
)

type EnhanceMode struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	state      EnhanceModeState
	tunDevice  *tun.TunDevice
	dnsServer  *dns.DNSServer
	dnsHandler *mockdns.DNSMock
}

func NewEnhanceMode() *EnhanceMode {
	return &EnhanceMode{state: EnhanceModeStateInit}
}

func (e *EnhanceMode) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	e.ctx = ctx
	e.cancel = cancel
	e.state = EnhanceModeStateRunning

	tunDevice, err := tun.OpenTun(ctx)
	if err != nil {
		return err
	}
	e.tunDevice = tunDevice

	// dns mock server
	e.dnsHandler = mockdns.NewMockHandle()
	e.dnsServer = dns.NewDNSServer(e.dnsHandler)

	go e.handleUdp()
	go e.handleTcp()

	return nil
}

func (e *EnhanceMode) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.state = EnhanceModeStateStopped
	if e.cancel != nil {
		e.cancel()
	}
	if e.tunDevice != nil {
		e.tunDevice.Close()
	}
	return nil
}
