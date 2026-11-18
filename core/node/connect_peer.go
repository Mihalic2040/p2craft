package node

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (n *Node) ConnectPeer(ctx context.Context, id peer.ID) error {
	// Already connected?
	if n.Host.Network().Connectedness(id) == network.Connected {
		log.Printf("[P2P] Already connected to %s", id)
		return nil
	}

	// Check whether we already know this peer's addresses.
	addrs := n.Host.Peerstore().Addrs(id)

	if len(addrs) == 0 {
		log.Printf("[DHT] Resolving peer %s", id)

		info, err := n.Idht.FindPeer(ctx, id)
		if err != nil {
			return fmt.Errorf("DHT FindPeer (%s): %w", id, err)
		}

		addrs = info.Addrs
	}

	if len(addrs) == 0 {
		return fmt.Errorf("no addresses found for peer %s", id)
	}

	log.Printf("[P2P] Found %d addresses for %s", len(addrs), id)

	info := peer.AddrInfo{
		ID:    id,
		Addrs: addrs,
	}

	if err := n.Host.Connect(ctx, info); err != nil {
		return fmt.Errorf("connecting to %s: %w", id, err)
	}

	log.Printf("[P2P] Successfully connected to %s", id)
	return nil
}
