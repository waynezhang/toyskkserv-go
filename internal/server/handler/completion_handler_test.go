package handler

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
)

func TestCompletionHandler(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	w := bytes.NewBuffer(nil)

	dm := dictionary.NewDictManager(tmp, false)
	dm.DictionariesDidChange([]string{
		"../../../testdata/jisyo.utf8",
		"../../../testdata/jisyo-2.utf8",
	})
	h := CompletionHandler{dm: dm}

	w.Reset()
	assert.True(t, h.Do("tai", w))
	assert.Equal(t, "1/taiwan/\n", w.String())

	w.Reset()
	assert.True(t, h.Do("1", w))
	assert.Equal(t, "1/1024/1234/1seg/\n", w.String())

	w.Reset()
	assert.True(t, h.Do("tawww", w))
	assert.Equal(t, "4tawww \n", w.String())
}
