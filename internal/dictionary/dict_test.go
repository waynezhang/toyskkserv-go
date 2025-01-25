package dictionary

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, ENCODING_UNDECIDED, parseEncoding([]byte("")))
	assert.Equal(t, ENCODING_UNDECIDED, parseEncoding([]byte("xxx")))
	assert.Equal(t, ENCODING_UNDECIDED, parseEncoding([]byte(";; -*- coding -*-")))
	assert.Equal(t, ENCODING_EUCJP, parseEncoding([]byte(";; -*- coding: euc-jis-2004 -*-")))
	assert.Equal(t, ENCODING_UTF8, parseEncoding([]byte(";; -*- coding: utf-8 -*-")))
}
