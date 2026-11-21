package node

import (
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/event"
)

func (n *Node) MonitorDHT() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-n.Ctx.Done():
			return

		case <-ticker.C:
			dhtPeers := n.Idht.RoutingTable().ListPeers()
			connectedPeers := n.Host.Network().Peers()

			log.Printf(
				"[DHT] Routing peers: %d | Connected peers: %d",
				len(dhtPeers),
				len(connectedPeers),
			)
		}
	}
}

func (n *Node) WatchAutoRelay() error {
	h := n.Host
	ctx := n.Ctx
	sub, err := h.EventBus().Subscribe(
		new(event.EvtAutoRelayAddrsUpdated),
	)
	if err != nil {
		return err
	}

	go func() {
		defer sub.Close()

		log.Println("[AUTORELAY] Waiting for relay reservations...")

		for {
			select {
			case <-ctx.Done():
				return

			case e, ok := <-sub.Out():
				if !ok {
					return
				}

				evt := e.(event.EvtAutoRelayAddrsUpdated)

				if len(evt.RelayAddrs) == 0 {
					log.Println("[AUTORELAY] No active relay reservations")
					continue
				}

				log.Printf(
					"[AUTORELAY] %d active relay address(es)",
					len(evt.RelayAddrs),
				)

				for _, addr := range evt.RelayAddrs {
					log.Printf("[AUTORELAY] Available: %s", addr)
				}
			}
		}
	}()

	return nil
}
