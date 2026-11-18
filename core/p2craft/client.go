package p2craft

import (
	"context"
	"log"
	"net"
)

// Start starts a local TCP listener.
func (p *TCPProxy) Start(addr string) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		p.cancel()
		return err
	}

	p.listener = listener

	log.Printf("[PROXY] Listening on %s", listener.Addr())

	p.wg.Add(1)
	go p.acceptLoop()

	return nil
}

func (p *TCPProxy) Addr() string {
	return p.listener.Addr().String()
}

func (p *TCPProxy) acceptLoop() {
	defer p.wg.Done()

	for {
		conn, err := p.listener.Accept()
		if err != nil {
			select {
			case <-p.ctx.Done():
				return
			default:
				log.Printf("[PROXY] Accept error: %v", err)
				continue
			}
		}

		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.handleConnection(conn)
		}()
	}
}

func (p *TCPProxy) handleConnection(conn net.Conn) {
	defer conn.Close()

	stream, err := p.host.NewStream(
		p.ctx,
		p.remote,
		MinecraftProtocol,
	)
	log.Println("remote", p.remote)
	if err != nil {
		log.Printf("[PROXY] Failed to open stream: %v", err)
		return
	}
	defer stream.Close()

	log.Printf("[PROXY] Connected to %s", p.remote)

	Forward(conn, stream)
}

func (p *TCPProxy) Close() error {
	var err error

	p.once.Do(func() {
		if p.cancel != nil {
			p.cancel()
		}

		if p.listener != nil {
			err = p.listener.Close()
		}

		p.wg.Wait()
	})

	return err
}
