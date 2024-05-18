package p2p

// Peerer is an interface that represents the remote node.
type Peerer interface{}

// Transporter is anything that handles communication between
// the nodes in the network. This can be of the form of TCP,
// UDP, websockets etc.
type Transporter interface {
	ListenAndAccept() error
}
