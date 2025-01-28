package handler

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
		"../../../testdata/jisyo.utf8",
		"../../../testdata/jisyo-2.utf8",
	})
	h := CandidateHandler{dm: dm}

	w.Reset()
	assert.True(t, h.Do("taiwan", w))
	assert.Equal(t, "1/台湾/🇹🇼/\n", w.String())

	w.Reset()
	assert.True(t, h.Do("tai", w))
	assert.Equal(t, "4tai \n", w.String())
}

// helper func
func prepareTempDir(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "toyskkserv-test")
	assert.Nil(t, err)

	return tmp
}
