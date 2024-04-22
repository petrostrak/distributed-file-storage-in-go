package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peerer
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: addr,
	}
}
