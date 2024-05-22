package main

import (
	"log"

	"github.com/petrostrak/distributed-file-storage-in-go/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOps{
		ListenAddr: ":3000",
		ShakeHands: p2p.NoHandshake,
		Decoder:    p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	s := NewFileServer(fileServerOpts)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
