package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOps{
		ListenAddr: ":3000",
		ShakeHands: NoHandshake,
		Decoder:    DefaultDecoder{},
	}

	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.TCPTransportOptions.ListenAddr, ":3000")

	// Server

	assert.Nil(t, tr.ListenAndAccept())
}
