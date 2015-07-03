package uats

import (
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestGetServer(t *testing.T) {
	client := clc.New(clc.EnvConfig())
	server, err := client.Server.Get("VA1T3BKAPI01")

	assert.Nil(t, err)
	assert.Equal(t, "VA1T3BKAPI01", server.Name)
}
