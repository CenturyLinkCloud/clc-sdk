package clc_test

import (
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestGetID_Success(t *testing.T) {
	assert := assert.New(t)

	links := clc.Links{clc.Link{Rel: "a", ID: "1"}}
	ok, id := links.GetID("a")

	assert.True(ok)
	assert.Equal(id, "1")
}

func TestGetID_Failure(t *testing.T) {
	assert := assert.New(t)

	links := clc.Links{clc.Link{Rel: "a", ID: "1"}}
	ok, id := links.GetID("b")

	assert.False(ok)
	assert.Equal(id, "")
}
