package server

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
)

func TestCandidateHandler(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	w := bytes.NewBuffer(nil)

	dm := dictionary.NewDictManager(tmp, false)
	dm.DictionariesDidChange([]string{
		"../../testdata/jisyo.utf8",
		"../../testdata/jisyo-2.utf8",
	})
	h := candidateHandler{dm: dm}

	w.Reset()
	assert.True(t, h.do("taiwan", w))
	assert.Equal(t, "1/台湾/🇹🇼/\n", w.String())

	w.Reset()
	assert.True(t, h.do("tai", w))
	assert.Equal(t, "4tai \n", w.String())
}
