package dictionary

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/tskks/internal/iconv"
)

func TestLoadDict(t *testing.T) {
	cm := newCandidatesManager()
	loadFile("../../testdata/jisyo.utf8", cm)
	loadFile("../../testdata/jisyo.euc-jp", cm)
	assert.Equal(t, "/キロ/", cm.findCandidates("1024"))
	assert.Equal(t, "/ā;a-/å;a^/ä;a:/ã;a~/â;a^/á;a'/à;a`/ă;av/ą;a,/ⓐ;(a)/ª;西語女性序数/ɐ;[IPA]/ʌ;[IPA]/ɑ;[IPA]/ɒ;[IPA]/", cm.findCandidates("a"))
}

func TestLoadInvalidDict(t *testing.T) {
	cm := newCandidatesManager()
	loadFile("../../testdata/jisyo.utf8.notexisting", cm)
	assert.Equal(t, "", cm.findCandidates("1024"))
}

func TestParseEncoding(t *testing.T) {
	assert.Equal(t, iconv.ENCODING_UNDECIDED, parseEncoding(""))
	assert.Equal(t, iconv.ENCODING_UNDECIDED, parseEncoding("xxx"))
	assert.Equal(t, iconv.ENCODING_UNDECIDED, parseEncoding(";; -*- coding -*-"))
	assert.Equal(t, iconv.ENCODING_EUCJP, parseEncoding(";; -*- coding: euc-jis-2004 -*-"))
	assert.Equal(t, iconv.ENCODING_UTF8, parseEncoding(";; -*- coding: utf-8 -*-"))
}
