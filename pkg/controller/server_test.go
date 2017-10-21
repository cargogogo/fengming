package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	assert := assert.New(t)

	cfg := &ServerConfig{
		Addr:             "0.0.0.0:9090",
		RegistryBlobPath: "/home/ubuntu/registry/docker/registry/v2/blobs/sha256",
		TrackerAddr:      "45.76.163.62:6881",
	}
	server, err := NewServer(cfg)
	assert.NoError(err)

	assert.NoError(server.Run())
}
