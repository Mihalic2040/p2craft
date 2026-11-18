package node

import (
	"log"
	"time"
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
