package server

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/toyskkserv/internal/defs"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
)

func TestNew(t *testing.T) {
	dm := dictionary.NewDictManager("tmp", false)
	s := New("addr", dm)

	assert.Equal(t, dm, s.dictManager)
	assert.Equal(t, "addr", s.listenAddr)
}

func TestHandleRequest(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	dm := dictionary.NewDictManager(tmp, false)
	dm.DictionariesDidChange([]string{
		"https://github.com/uasi/skk-emoji-jisyo/raw/refs/heads/master/SKK-JISYO.emoji.utf8",
	})

	cases := [][]interface{}{
		{"", true, ""},
		{"", true, " "},
		{"", false, "0"},
		{"1/😄/\n", true, "1smile "},
		{"1/zombie/zombie_man/zombie_woman/\n", true, "4zom "},
		{defs.VersionString() + " \n", true, "2 "},
		{"localhost:port \n", true, "3 "},
	}

	s := New("localhost:port", dm)
	w := bytes.NewBuffer(nil)
	for i, c := range cases {
		msg := "case " + strconv.Itoa(i)
		w.Reset()

		running := s.handleRequest(c[2].(string), w)
		resp := w.String()

		assert.Equal(t, c[0], resp, msg)
		assert.Equal(t, c[1], running, msg)
	}
}

// helper func
func prepareTempDir(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "toyskkserv-test")
	assert.Nil(t, err)

	return tmp
}
