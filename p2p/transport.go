package p2p

import "net"

// Peerer is an interface that represents the remote node.
type Peerer interface {
	RemoteAddr() net.Addr
	Close() error
}

// Transporter is anything that handles communication between
// the nodes in the network. This can be of the form of TCP,
// UDP, websockets etc.
type Transporter interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
