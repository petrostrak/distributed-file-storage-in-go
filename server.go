package main

import "github.com/petrostrak/distributed-file-storage-in-go/p2p"

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
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		RootDir:           opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
	}
}

func (s *FileServerOpts) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}
