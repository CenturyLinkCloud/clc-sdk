package clc_test

import (
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestInitializeClient(t *testing.T) {
	assert := assert.New(t)

	client := clc.New(clc.Config{})

	assert.NotNil(client, "expected client to not be nil")
}

func TestAuth(t *testing.T) {
	assert := assert.New(t)

	client := clc.New(clc.Config{})
	_, err := client.Auth()

	assert.Nil(err, "no error expected during auth")
}
