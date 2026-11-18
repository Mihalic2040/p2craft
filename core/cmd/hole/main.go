package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"demo/node"
	"demo/p2craft"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	node := node.NewNode(true, os.Getenv("HOLE_PRIV_KEY"))
	if err := node.Start(); err != nil {
		panic(err)
	}

	hole_id := node.Host.ID().String() + ".p2p"
	log.Println("Hole ID:", hole_id)

	go node.MonitorDHT()
	p2craft.RegisterMinecraftHandler(node.Host, "0.0.0.0:25565")
	log.Println("Service started.")

	<-ctx.Done()
	log.Println("Shutdown signal received.")

	// Stop accepting new connections and close existing ones.
	if err := node.Host.Close(); err != nil {
		log.Printf("Error closing libp2p host: %v", err)
	}
}
