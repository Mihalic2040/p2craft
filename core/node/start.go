package node

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

func (n *Node) Start() error {
	// New ctx new life
	n.Ctx = context.Background()

	if err := n.Idht.Bootstrap(n.Ctx); err != nil {
		n.Host.Close()
		log.Println("Fail Bootstrap DHT:", err)
		return err
	}

	// Local network discovery
	discovery := mdns.NewMdnsService(
		n.Host,
		"p2craft",
		&discoveryNotifee{h: n.Host},
	)

	if err := discovery.Start(); err != nil {
		n.Host.Close()
		log.Println("Fail to start mdns discovery:", err)
		return err
	}

	log.Println("MDNS Started.")

	go RunAutoRelayFeeder(
		n.Ctx,
		n.Host,
		n.Idht,
		nil, // trusted peers
		n.peer_cnah,
	)

	log.Println("RelayFeeder Started.")

	return nil
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("MDNS Discovered peer: %s", pi.ID)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := n.h.Connect(ctx, pi); err != nil {
		log.Printf("Failed to connect to peer %s: %v", pi.ID, err)
	}
}
