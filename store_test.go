package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsPics"
	pathname := CASPathTransformFunc(key)
	expected := "5f30a/6b2be/aff4a/6a4ee/f5506/0bd74/6444e/a54c6"

	if pathname != expected {
		t.Errorf("have %s want %s\n", pathname, expected)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}

	s := NewStore(opts)

	data := bytes.NewReader([]byte("somejpgbytes"))
	if err := s.writeStream("mySpecialPic", data); err != nil {
		t.Error(err)
	}
}
