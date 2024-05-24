package p2p

type HandshakeFunc func(Peer) error

func NoHandshake(peer Peer) error { return nil }
