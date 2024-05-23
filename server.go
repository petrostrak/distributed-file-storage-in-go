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
	Transport      p2p.Transporter
	BootstrapNodes []string
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

	s.bootstrapNetwork()
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

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error", err)
			}
		}(addr)
	}

	return nil
}
