package handler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostHandler(t *testing.T) {
	w := bytes.NewBuffer(nil)
	assert.True(t, HostHandler{host: "127.0.0.1:9999"}.Do("", w))
	assert.Equal(t, "127.0.0.1:9999 \n", w.String())
}
