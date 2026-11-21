package node

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

type Node struct {
	Server bool

	Ctx        context.Context
	privateKey crypto.PrivKey

	peer_cnah chan<- peer.AddrInfo

	Host host.Host
	Idht *dht.IpfsDHT
}

func NewNode(server bool, priv_key_path string) *Node {
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	if server {
		priv, err = LoadOrCreateIdentity(priv_key_path)
		if err != nil {
			panic(err)
		}
	}

	connmgr, err := connmgr.NewConnManager(
		50,  // Lowwater
		150, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	var idht *dht.IpfsDHT
	var dhtMode dht.ModeOpt = dht.ModeAuto
	if !server {
		dhtMode = dht.ModeClient
	}

	peerSource, peerChan := NewAutoRelayPeerSource()

	h, err := libp2p.New(
		// Persistent identity
		libp2p.Identity(priv),

		// Only native transports
		libp2p.NoTransports,
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(quic.NewTransport),
		libp2p.DefaultListenAddrs,

		// Noise encryption
		libp2p.Security(noise.ID, noise.New),

		// Connection management
		libp2p.ConnectionManager(connmgr),

		// NAT traversal
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
		libp2p.EnableAutoNATv2(),
		// libp2p.EnableAutoRelayWithPoeerSource(),
		libp2p.EnableRelay(),
		libp2p.EnableAutoRelayWithPeerSource(
			peerSource,
			autorelay.WithNumRelays(1),
			autorelay.WithMinCandidates(2),
		),
		libp2p.EnableHolePunching(),

		// Small DHT: client mode
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			idht, err = dht.New(
				h,
				dht.Mode(dhtMode),
				dht.BootstrapPeers(dht.GetDefaultBootstrapPeerAddrInfos()...),
			)
			return idht, err
		}),
	)
	if err != nil {
		panic(err)
	}

	node := Node{
		privateKey: priv,
		Host:       h,
		Idht:       idht,
		peer_cnah:  peerChan,
	}

	// The last step to get fully up and running would be to connect to
	// bootstrap peers (or any other peers). We leave this commented as
	// this is an example and the peer will die as soon as it finishes, so
	// it is unnecessary to put strain on the network.

	/*
		// This connects to public bootstrappers
		for _, addr := range dht.DefaultBootstrapPeers {
			pi, _ := peer.AddrInfoFromP2pAddr(addr)
			// We ignore errors as some bootstrap peers may be down
			// and that is fine.
			h2.Connect(ctx, *pi)
		}
	*/
	return &node
}
