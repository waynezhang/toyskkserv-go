package dictionary

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/toyskkserv/internal/dictionary/candidate"
)

func TestLoadDict(t *testing.T) {
	m := candidate.New(false)
	loadFile("../../testdata/jisyo.utf8", m)
	loadFile("../../testdata/jisyo.euc-jp", m)
	assert.Equal(t, "/̀;accent grave (diacritic)/", m.Find("`")) // first line
	assert.Equal(t, "/キロ/", m.Find("1024"))
	assert.Equal(t, "/ā;a-/å;a^/ä;a:/ã;a~/â;a^/á;a'/à;a`/ă;av/ą;a,/ⓐ;(a)/ª;西語女性序数/ɐ;[IPA]/ʌ;[IPA]/ɑ;[IPA]/ɒ;[IPA]/", m.Find("a"))
}

func TestLoadInvalidDict(t *testing.T) {
	m := candidate.New(false)
	loadFile("../../testdata/jisyo.utf8.notexisting", m)
	assert.Equal(t, "", m.Find("1024"))
}

func TestDetectEncoding(t *testing.T) {
	cases := [][]string{
		{"jisyo.euc-jp", ENCODING_EUCJP},
		{"jisyo.euc-jp.withheader", ENCODING_EUCJP},
		{"jisyo.euc-jp.empty", ENCODING_EUCJP},
		{"jisyo.utf8.withheader", ENCODING_UTF8},
	}
	for _, c := range cases {
		f, err := os.Open("../../testdata/" + c[0])
		assert.Nil(t, err)
		defer f.Close()

		enc, err := detectFileEncoding(f)
		assert.Nil(t, err)
		assert.Equal(t, c[1], enc)

		pos, err := f.Seek(0, io.SeekCurrent)
		assert.Equal(t, int64(0), pos)
	}
}
