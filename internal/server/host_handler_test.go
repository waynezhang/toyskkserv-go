package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostHandler(t *testing.T) {
	w := bytes.NewBuffer(nil)
	assert.True(t, hostHandler{host: "127.0.0.1:9999"}.do("", w))
	assert.Equal(t, "127.0.0.1:9999 \n", w.String())
}
