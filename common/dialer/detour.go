package dialer

import (
	"context"
	"net"

	"github.com/sagernet/sing-box/adapter"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

type DetourDialer struct {
	router adapter.Router
	detour string
}

func NewDetour(router adapter.Router, detour string) N.Dialer {
	return &DetourDialer{router: router, detour: detour}
}

func (d *DetourDialer) Start() error {
	_, err := d.Dialer()
	return err
}

func (d *DetourDialer) Dialer() (N.Dialer, error) {
	var err error
	detour, loaded := d.router.Outbound(d.detour)
	if !loaded {
		err = E.New("outbound detour not found: ", d.detour)
	}
	return detour, err
}

func (d *DetourDialer) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	dialer, err := d.Dialer()
	if err != nil {
		return nil, err
	}
	return dialer.DialContext(ctx, network, destination)
}

func (d *DetourDialer) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	dialer, err := d.Dialer()
	if err != nil {
		return nil, err
	}
	return dialer.ListenPacket(ctx, destination)
}

func (d *DetourDialer) Upstream() any {
	detour, _ := d.Dialer()
	return detour
}
