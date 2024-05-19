package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer.
	conn net.Conn

	// if we dial and retrieve a conn => outbound == true.
	// if we accept and retrieve a conn => outbound == false.
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbount bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbount,
	}
}

type TCPTransportOps struct {
	ListenAddr string
	ShakeHands HandshakeFunc
	Decoder    Decoder
}

type TCPTransport struct {
	TCPTransportOptions TCPTransportOps
	listener            net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peerer
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.TCPTransportOptions.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %v", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.TCPTransportOptions.ShakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("tcp handshake error: %s\v", err)
		return
	}

	// Read loop
	msg := &RPC{}
	for {
		if err := t.TCPTransportOptions.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("tcp error: %s\n", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("message: %+v\n", msg)
	}

}
