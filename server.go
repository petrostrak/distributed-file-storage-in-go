package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

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

	mu    sync.Mutex
	peers map[string]p2p.Peerer
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
		peers:          make(map[string]p2p.Peerer),
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
			var data Message
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&data); err != nil {
				log.Fatal(err)
			}

			if err := s.handleMessage(&data); err != nil {
				log.Println(err)
			}
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

func (s *FileServer) OnPeer(p p2p.Peerer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s\n", p.RemoteAddr())

	return nil
}

type Message struct {
	From    string
	Payload any
}

type DataMessage struct {
	Key  string
	Data []byte
}

func (s *FileServer) StoreFile(key string, r io.Reader) error {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)

	// 1. Store the file to disk
	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	p := &DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}

	// 2. broadcast the file to all known peers in the network
	return s.broadcast(&Message{
		From:    "todo",
		Payload: p,
	})
}

func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) handleMessage(msg *Message) error {
	switch v := msg.Payload.(type) {
	case *DataMessage:
		fmt.Printf("received data %+v\n", v)
	}
	return nil
}
