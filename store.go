package main

import (
	"io"
	"log"
	"os"
)

var DefaultPathTransformFunc = func(key string) string { return key }

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	if err := os.MkdirAll(key, os.ModePerm); err != nil {
		return err
	}

	filename := "somefilename"
	filename = pathName + filename

	f, err := os.Create(pathName + "/" + filename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s\n", n, filename)

	return nil
}
