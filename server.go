package main

import (
	"fmt"
	"log"

	"github.com/petrostrak/distributed-file-storage-in-go/p2p"
)

type FileServerOpts struct {
	ListenAddr  string
	StorageRoot string
	PathTransformFunc
	Transport        p2p.Transporter
	TCPTransportOpts p2p.TCPTransportOps
}

type FileServer struct {
	FileServerOpts
	store *Store
	quit  chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		RootDir:           opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quit:           make(chan struct{}),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quit)
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("fileserver stopped due to quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quit:
			return
		}
	}
}
