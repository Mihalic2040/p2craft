package node

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
)

// ClosestPeerFinder is implemented by a Kademlia DHT.
type ClosestPeerFinder interface {
	GetClosestPeers(ctx context.Context, key string) ([]peer.ID, error)
}

// NewAutoRelayPeerSource creates the candidate source consumed by AutoRelay
// and a channel used by the background feeder.
func NewAutoRelayPeerSource() (
	autorelay.PeerSource,
	chan<- peer.AddrInfo,
) {
	peerChan := make(chan peer.AddrInfo)

	source := func(ctx context.Context, num int) <-chan peer.AddrInfo {
		out := make(chan peer.AddrInfo)

		go func() {
			defer close(out)

			for i := 0; i < num; i++ {
				select {
				case p, ok := <-peerChan:
					if !ok {
						return
					}

					select {
					case out <- p:
					case <-ctx.Done():
						return
					}

				case <-ctx.Done():
					return
				}
			}
		}()

		return out
	}

	return source, peerChan
}

// RunAutoRelayFeeder periodically feeds potential relay candidates
// into the AutoRelay PeerSource.
//
// Candidate sources:
//   - Trusted peers
//   - DHT closest peers
//   - Currently connected peers
func RunAutoRelayFeeder(
	ctx context.Context,
	h host.Host,
	dht ClosestPeerFinder,
	trustedPeers []peer.AddrInfo,
	peerChan chan<- peer.AddrInfo,
) {
	defer close(peerChan)

	delay := 15 * time.Second

	for {
		timer := time.NewTimer(delay)

		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return
		}

		// 1. Feed trusted peers.
		for _, p := range trustedPeers {
			if len(p.Addrs) == 0 {
				continue
			}

			select {
			case peerChan <- p:
			case <-ctx.Done():
				return
			}
		}

		// 2. Feed peers discovered through DHT.
		if dht != nil {
			closestPeers, err := dht.GetClosestPeers(
				ctx,
				h.ID().String(),
			)

			if err == nil {
				for _, id := range closestPeers {
					addrs := h.Peerstore().Addrs(id)
					if len(addrs) == 0 {
						continue
					}

					p := peer.AddrInfo{
						ID:    id,
						Addrs: addrs,
					}

					select {
					case peerChan <- p:
					case <-ctx.Done():
						return
					}
				}
			} else if ctx.Err() == nil {
				log.Printf("AutoRelay: DHT discovery failed: %v", err)
			}
		}

		// 3. Feed all currently connected peers.
		for _, id := range h.Network().Peers() {
			addrs := h.Peerstore().Addrs(id)
			if len(addrs) == 0 {
				continue
			}

			p := peer.AddrInfo{
				ID:    id,
				Addrs: addrs,
			}

			select {
			case peerChan <- p:
			case <-ctx.Done():
				return
			}
		}

		// Exponential backoff: 15s -> 45s -> 135s -> ... -> 1h.
		delay *= 3
		if delay > time.Hour {
			delay = time.Hour
		}
	}
}
