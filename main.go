package main

import (
	"fmt"
	"log"

	"github.com/petrostrak/distributed-file-storage-in-go/p2p"
)

func main() {
	opts := p2p.TCPTransportOps{
		ListenAddr: ":3000",
		ShakeHands: p2p.NoHandshake,
		Decoder:    p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(opts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
