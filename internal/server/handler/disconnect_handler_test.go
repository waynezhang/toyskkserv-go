package handler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisconnectHandler(t *testing.T) {
	assert.False(t, DisconnectHandler{}.Do("", bytes.NewBuffer(nil)))
}
