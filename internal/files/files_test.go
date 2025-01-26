package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileChecksum(t *testing.T) {
	checksum, err := FileChecksum("../../testdata/checksum_testfile")

	assert.Nil(t, err)
	assert.Equal(t,
		"e0ef7229e64c61596d8be928397e19fcc542ac920c4132106fb1ec2295dd73d1",
		checksum)

	checksum, err = FileChecksum("../../testdata/checksum_testfile_not_existing")
	assert.NotNil(t, err)
	assert.Equal(t, "", "")
}

func TestIsFileExisting(t *testing.T) {
	assert.True(t, IsFileExisting("../../testdata/checksum_testfile"))
	assert.False(t, IsFileExisting("../../testdata/checksum_testfile_not_existing"))
}

func TestDictName(t *testing.T) {
	assert.Equal(
		t,
		"SKK-JISYO.lisp",
		DictName("https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.lisp"),
	)
}

func TestDictionaryPaths(t *testing.T) {
	urls := []string{
		"https://test.com/dict1.abc",
		"https://test.com/dict2.def",
		"/usr/bin",
		"https://test.com/dict3.ghi",
		"/usr/local/bin",
		"https://test.com/dict4.jkl",
		"~/abc",
	}

	home, err := os.UserHomeDir()
	assert.Nil(t, err)

	paths := DictionaryPaths(urls, "/dir")
	assert.Equal(t, []string{
		"/dir/dict1.abc",
		"/dir/dict2.def",
		"/usr/bin",
		"/dir/dict3.ghi",
		"/usr/local/bin",
		"/dir/dict4.jkl",
		home + "/abc",
	}, paths)
}

func TestIsLocalURL(t *testing.T) {
	assert.True(t, IsLocalURL("/usr/local/bin"))
	assert.True(t, IsLocalURL("../../testdata"))
	assert.True(t, IsLocalURL("~/.config"))
	assert.False(t, IsLocalURL("http://github.com"))
	assert.False(t, IsLocalURL("https://github.com"))
}
