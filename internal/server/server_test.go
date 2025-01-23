package server

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/tskks/internal/defs"
	"github.com/waynezhang/tskks/internal/dictionary"
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
		{"4/ \n", true, "1"},
		{"4/1 \n", true, "11"},
		{"4/1 \n", true, "11 "},
		{"1/😄/\n", true, "1smile "},
		{"1/zombie/zombie_man/zombie_woman/\n", true, "4zom "},
		{"4/somethingnotexisted \n", true, "4somethingnotexisted "},
		{defs.VersionString() + " \n", true, "2"},
		{defs.VersionString() + " \n", true, "2 "},
		{"localhost:port \n", true, "3"},
		{"localhost:port \n", true, "3 "},
	}

	s := New("localhost:port", dm)
	for i, c := range cases {
		msg := "case " + strconv.Itoa(i)
		resp, running := s.handleRequest(c[2].(string))
		assert.Equal(t, c[0], resp, msg)
		assert.Equal(t, c[1], running, msg)
	}
}

// helper func
func prepareTempDir(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "tskks-test")
	assert.Nil(t, err)

	return tmp
}
