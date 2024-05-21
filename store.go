package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

var DefaultPathTransformFunc = func(key string) string { return key }

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc
}

type Store struct {
	StoreOpts
}

func CASPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")
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
