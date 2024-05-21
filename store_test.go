package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsPics"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "5f30a6b2beaff4a6a4eef55060bd746444ea54c6"
	expectedPathname := "5f30a/6b2be/aff4a/6a4ee/f5506/0bd74/6444e/a54c6"

	if pathKey.Pathname != expectedPathname {
		t.Errorf("have %s want %s\n", pathKey.Pathname, expectedPathname)
	}

	if pathKey.Original != expectedOriginalKey {
		t.Errorf("have %s want %s\n", pathKey.Original, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)

	data := bytes.NewReader([]byte("somejpgbytes"))
	if err := s.writeStream("mySpecialPic", data); err != nil {
		t.Error(err)
	}
}
