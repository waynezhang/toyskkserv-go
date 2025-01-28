package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/toyskkserv/internal/defs"
)

func TestVersionHandler(t *testing.T) {
	w := bytes.NewBuffer(nil)
	assert.True(t, versionHandler{}.do("", w))
	assert.Equal(t, defs.VersionString()+" \n", w.String())
}
