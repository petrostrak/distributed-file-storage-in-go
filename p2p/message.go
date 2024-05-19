package p2p

// Message holds any arbitrary data that are sent over each
// each transport between two nodes in the network.
type Message struct {
	Payload []byte
}
