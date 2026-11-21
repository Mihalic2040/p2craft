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

	target_addres := os.Getenv("HOLE_TARGET_ADDR")
	if target_addres == "" {
		log.Println("[WARN]", "No target addr specified. Using default 127.0.0.1:25565")
		target_addres = "127.0.0.1:25565"
	}

	hole_id := node.Host.ID().String() + ".p2p"
	log.Println("Hole ID:", hole_id)

	for _, addr := range node.Host.Addrs() {
		log.Println("Listening on:", addr)
	}

	go node.MonitorDHT()
	go node.WatchAutoRelay()
	p2craft.RegisterMinecraftHandler(node.Host, target_addres)
	log.Println("Service started.")

	<-ctx.Done()
	log.Println("Shutdown signal received.")

	// Stop accepting new connections and close existing ones.
	if err := node.Host.Close(); err != nil {
		log.Printf("Error closing libp2p host: %v", err)
	}
}
