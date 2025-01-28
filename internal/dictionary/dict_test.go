package dictionary

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDict(t *testing.T) {
	cm := newCandidatesManager()
	loadFile("../../testdata/jisyo.utf8", cm)
	loadFile("../../testdata/jisyo.euc-jp", cm)
	assert.Equal(t, "/̀;accent grave (diacritic)/", cm.findCandidates("`")) // first line
	assert.Equal(t, "/キロ/", cm.findCandidates("1024"))
	assert.Equal(t, "/ā;a-/å;a^/ä;a:/ã;a~/â;a^/á;a'/à;a`/ă;av/ą;a,/ⓐ;(a)/ª;西語女性序数/ɐ;[IPA]/ʌ;[IPA]/ɑ;[IPA]/ɒ;[IPA]/", cm.findCandidates("a"))
}

func TestLoadInvalidDict(t *testing.T) {
	cm := newCandidatesManager()
	loadFile("../../testdata/jisyo.utf8.notexisting", cm)
	assert.Equal(t, "", cm.findCandidates("1024"))
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
