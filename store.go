package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	defaultRootFolderName = "network"
)

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) fullpath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	RootDir string
	PathTransformFunc
}

type Store struct {
	StoreOpts
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.RootDir) == 0 {
		opts.RootDir = defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(s.RootDir+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}

	filename := pathKey.fullpath()

	f, err := os.Create(s.RootDir + "/" + filename)
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

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(s.RootDir + "/" + pathKey.fullpath())
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (p PathKey) rootDir() string {
	path := strings.Split(p.Pathname, "/")
	if len(path) == 0 {
		return ""
	}

	return path[0]
}

func (s *Store) Delete(key string) error {
	pathkey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk\n", pathkey.Filename)
	}()

	path := fmt.Sprintf("%s/%s", s.RootDir, pathkey.rootDir())

	return os.RemoveAll(path)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.RootDir)
}

func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformFunc(key)
	_, err := os.Stat(s.RootDir + "/" + pathkey.fullpath())

	return !errors.Is(err, os.ErrNotExist)
}
