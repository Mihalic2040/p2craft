package p2craft

import (
	"context"
	"net"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type TCPProxy struct {
	host   host.Host
	remote peer.ID

	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc

	wg   sync.WaitGroup
	once sync.Once
}

func NewTCPProxy(h host.Host, remote peer.ID) *TCPProxy {
	return &TCPProxy{
		host:   h,
		remote: remote,
	}
}
