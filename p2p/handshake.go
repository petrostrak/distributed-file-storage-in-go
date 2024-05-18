package p2p

type HandshakeFunc func(Peerer) error

func NoHandshake(peer Peerer) error { return nil }
