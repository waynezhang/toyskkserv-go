package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisconnectHandler(t *testing.T) {
	assert.False(t, disconnectHandler{}.do("", bytes.NewBuffer(nil)))
}
