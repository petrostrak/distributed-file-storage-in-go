package main

import (
	"bytes"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "momsPics"
	pathKey := CASPathTransformFunc(key)
	expectedFilenameKey := "5f30a6b2beaff4a6a4eef55060bd746444ea54c6"
	expectedPathname := "5f30a/6b2be/aff4a/6a4ee/f5506/0bd74/6444e/a54c6"

	if pathKey.Pathname != expectedPathname {
		t.Errorf("have %s want %s\n", pathKey.Pathname, expectedPathname)
	}

	if pathKey.Filename != expectedFilenameKey {
		t.Errorf("have %s want %s\n", pathKey.Filename, expectedFilenameKey)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	key := "mySpecialPic"

	data := []byte("somejpgbytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s\n", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); ok {
		t.Errorf("expected NOT to have key %s\n", key)
	}
}
