package main

import (
	"bytes"
	"testing"
)

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
