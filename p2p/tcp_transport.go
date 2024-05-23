package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
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

// CLose implements the Peerer interface.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

type TCPTransportOps struct {
	ListenAddr string
	ShakeHands HandshakeFunc
	Decoder    Decoder
}

type TCPTransport struct {
	TCPTransportOptions TCPTransportOps
	listener            net.Listener
	rpcChan             chan RPC
	OnPeer              func(Peerer) error
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		rpcChan:             make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.TCPTransportOptions.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("tcp transport listening on port: %s\n", t.TCPTransportOptions.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error: %v", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.TCPTransportOptions.ShakeHands(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
		if err = t.TCPTransportOptions.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("tcp error: %s\n", err)
			return
		}

		rpc.From = conn.RemoteAddr()

		t.rpcChan <- rpc
	}

}

// Consume implements the Transport interface, which will return read-only channel
// for reading the incoming messages received from another peer in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChan
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the transport interface.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}
