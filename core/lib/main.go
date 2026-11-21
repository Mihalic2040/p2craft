package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"log"
	"sync"
	"unsafe"

	"demo/node"
	"demo/p2craft"

	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	StatusDisconnected = 0
	StatusConnected    = 1
	StatusError        = 2
)

type Tunnel struct {
	Proxy *p2craft.TCPProxy
	Addr  string
}

var (
	mu       sync.Mutex
	hostNode *node.Node
	tunnels  = make(map[string]*Tunnel)
	status   = StatusDisconnected
)

//export Connect
func Connect() C.int {
	mu.Lock()
	defer mu.Unlock()

	if status == StatusConnected {
		return 0
	}

	n := node.NewNode(false, "")

	go n.MonitorDHT()
	go n.WatchAutoRelay()

	if err := n.Start(); err != nil {
		log.Println(err)
		status = StatusError
		return -1
	}

	hostNode = n
	status = StatusConnected

	return 0
}

//export GetStatus
func GetStatus() C.int {
	mu.Lock()
	defer mu.Unlock()

	return C.int(status)
}

//export CreateTunnel
func CreateTunnel(remoteID *C.char) *C.char {
	id := C.GoString(remoteID)

	mu.Lock()
	defer mu.Unlock()

	if status != StatusConnected {
		return nil
	}

	if _, exists := tunnels[id]; exists {
		return C.CString(tunnels[id].Addr)
	}

	serverID, err := peer.Decode(id)
	if err != nil {
		return nil
	}

	proxy := p2craft.NewTCPProxy(hostNode.Host, hostNode.Idht, serverID)

	if err := proxy.Start("127.0.0.1:0"); err != nil {
		return nil
	}

	addr := proxy.Addr()

	tunnels[id] = &Tunnel{
		Proxy: proxy,
		Addr:  addr,
	}

	return C.CString(addr)
}

//export FreeString
func FreeString(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

//export Disconnect
func Disconnect() C.int {
	mu.Lock()
	defer mu.Unlock()

	for _, tunnel := range tunnels {
		tunnel.Proxy.Close()
	}

	tunnels = make(map[string]*Tunnel)

	if hostNode != nil {
		hostNode = nil
		for i := range tunnels {
			tunnel := tunnels[i]
			if err := tunnel.Proxy.Close(); err != nil {
				log.Println("Failed to close the tunnel for:", tunnel.Addr)
			}
		}
	}

	status = StatusDisconnected
	return 0
}

func main() {}
