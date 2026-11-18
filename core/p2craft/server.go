package p2craft

import (
	"log"
	"net"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

func RegisterMinecraftHandler(h host.Host, target string) {
	h.SetStreamHandler(MinecraftProtocol, func(s network.Stream) {
		defer s.Close()

		log.Printf("[PROXY] Incoming peer: %s", s.Conn().RemotePeer())
		log.Printf("[PROXY] Dialing target: %s", target)

		conn, err := net.DialTimeout("tcp", target, 5*time.Second)
		if err != nil {
			log.Printf("[PROXY] Target unavailable: %v", err)
			s.Reset()
			return
		}
		defer conn.Close()

		log.Printf("[PROXY] Target connected: %s", target)

		Forward(s, conn)

		log.Println("[PROXY] Connection closed")
	})
}
