package dictionary

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/tskks/internal/files"
)

func TestNew(t *testing.T) {
	dm := NewDictManager("tmp1", false)
	assert.NotNil(t, dm.cm)
	assert.Equal(t, "tmp1", dm.directory)
	assert.False(t, dm.fallbackToGoogle)

	dm = NewDictManager("tmp2", true)
	assert.NotNil(t, dm.cm)
	assert.Equal(t, "tmp2", dm.directory)
	assert.True(t, dm.fallbackToGoogle)
}

func TestHandleRequest(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	dm := NewDictManager(tmp, false)

	assert.Equal(t, "4/ ", dm.HandleRequest(""))
	assert.Equal(t, "4/ ", dm.HandleRequest(" "))
	assert.Equal(t, "4/ ", dm.HandleRequest("1"))
	assert.Equal(t, "4/abc ", dm.HandleRequest("1abc"))
	assert.Equal(t, "4/abc ", dm.HandleRequest("1abc "))

	dm.cm.addCandidates("abc", "/test1/test2/test3/")
	assert.Equal(t, "1/test1/test2/test3/", dm.HandleRequest("1abc"))

	dm.fallbackToGoogle = true
	assert.Equal(t, "1/アイウエオ/ア・イ・ウ・エ・オ/愛飢え男/aiueo/", dm.HandleRequest("1あいうえお"))
}

func TestHandleCompletion(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	dm := NewDictManager(tmp, false)

	assert.Equal(t, "4/ ", dm.HandleCompletion(""))
	assert.Equal(t, "4/ ", dm.HandleCompletion(" "))
	assert.Equal(t, "4/ ", dm.HandleCompletion("4"))
	assert.Equal(t, "4/abc ", dm.HandleCompletion("4abc"))
	assert.Equal(t, "4/abc ", dm.HandleCompletion("4abc "))

	dm.cm.addCandidates("abc", "/test1/test2/test3/")
	assert.Equal(t, "1/abc/", dm.HandleCompletion("4ab"))

	dm.cm.addCandidates("abd", "/test1/test2/test3/")
	assert.Equal(t, "1/abc/abd/", dm.HandleCompletion("4ab"))
}

func TestLoadAll(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	dm := NewDictManager(tmp, false)
	dm.reloadDicts([]string{
		"https://github.com/uasi/skk-emoji-jisyo/raw/refs/heads/master/SKK-JISYO.emoji.utf8",
	})

	assert.Equal(t, "1/👍/", dm.HandleRequest("1+1"))
}

func TestDownloadDictionary(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	dm := NewDictManager(tmp, false)
	dm.downloadDictionaries([]string{
		"https://github.com/uasi/skk-emoji-jisyo/raw/refs/heads/master/SKK-JISYO.emoji.utf8",
	})

	assert.True(t, files.IsFileExisting(filepath.Join(tmp, "SKK-JISYO.emoji.utf8")))
}

func TestLoadDictionaries(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	dm := NewDictManager(tmp, false)
	dm.loadDictionaries([]string{
		"../../testdata/jisyo.utf8",
		"../../testdata/jisyo.euc-jp",
	})

	cases := [][]string{
		{"/キロ/", "1024"},
		{"/ā;a-/å;a^/ä;a:/ã;a~/â;a^/á;a'/à;a`/ă;av/ą;a,/ⓐ;(a)/ª;西語女性序数/ɐ;[IPA]/ʌ;[IPA]/ɑ;[IPA]/ɒ;[IPA]/", "a"},
	}
	for idx, c := range cases {
		msg := "case " + strconv.Itoa(idx)
		cdd := dm.cm.findCandidates(c[1])
		assert.Equal(t, c[0], cdd, msg)
	}
}

func TestReloadDicts(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	dm := NewDictManager(tmp, false)
	dm.reloadDicts([]string{
		"https://github.com/uasi/skk-emoji-jisyo/raw/refs/heads/master/SKK-JISYO.emoji.utf8",
	})

	cases := [][]string{
		{"/🇹🇼/", "taiwan"},
	}
	for idx, c := range cases {
		msg := "case " + strconv.Itoa(idx)
		cdd := dm.cm.findCandidates(c[1])
		assert.Equal(t, c[0], cdd, msg)
	}

	dm.reloadDicts([]string{
		"https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan",
	})

	cases = [][]string{
		{"/台湾/", "taiwan"},
	}
	for idx, c := range cases {
		msg := "case " + strconv.Itoa(idx)
		cdd := dm.cm.findCandidates(c[1])
		assert.Equal(t, c[0], cdd, msg)
	}

	dm.reloadDicts([]string{
		"https://github.com/uasi/skk-emoji-jisyo/raw/refs/heads/master/SKK-JISYO.emoji.utf8",
		"https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan",
	})

	cases = [][]string{
		{"/🇹🇼/台湾/", "taiwan"},
	}
	for idx, c := range cases {
		msg := "case " + strconv.Itoa(idx)
		cdd := dm.cm.findCandidates(c[1])
		assert.Equal(t, c[0], cdd, msg)
	}

	dm.DictionariesDidChange([]string{
		"https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan",
	})
	cases = [][]string{
		{"/台湾/", "taiwan"},
	}
	for idx, c := range cases {
		msg := "case " + strconv.Itoa(idx)
		cdd := dm.cm.findCandidates(c[1])
		assert.Equal(t, c[0], cdd, msg)
	}
}

func TestDictionaryPaths(t *testing.T) {
	urls := []string{
		"https://test.com/dict1.abc",
		"https://test.com/dict2.def",
	}
	paths := dictionaryPaths(urls, "/dir")
	assert.Equal(t, []string{"/dir/dict1.abc", "/dir/dict2.def"}, paths)
}

// helper func
func prepareTempDir(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "tskks-test")
	assert.Nil(t, err)

	return tmp
}
