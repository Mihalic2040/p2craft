package main

import (
	"fmt"
	"log"
	"os"

	"demo/node"
	"demo/p2craft"

	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	node := node.NewNode(false, "")
	if err := node.Start(); err != nil {
		panic(err)
	}
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <name>")
		os.Exit(1)
	}

	serverID, err := peer.Decode(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	go node.MonitorDHT()
	go node.WatchAutoRelay()

	tcp_hole := p2craft.NewTCPProxy(node.Host, node.Idht, serverID)
	tcp_hole.Start("127.0.0.1:6666")

	select {}
}
