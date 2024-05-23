package p2p

// Peerer is an interface that represents the remote node.
type Peerer interface {
	Close() error
}

// Transporter is anything that handles communication between
// the nodes in the network. This can be of the form of TCP,
// UDP, websockets etc.
type Transporter interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
