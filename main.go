package main

import (
	"log"

	"github.com/petrostrak/distributed-file-storage-in-go/p2p"
)

func main() {
	opts := p2p.TCPTransportOps{
		ListenAddr: ":3000",
		ShakeHands: p2p.NoHandshake,
		Decoder:    p2p.GOBDecoder{},
	}

	tr := p2p.NewTCPTransport(opts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
