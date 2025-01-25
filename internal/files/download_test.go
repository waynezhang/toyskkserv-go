package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	file := filepath.Join(tmp, "file1")
	u, err := Download("https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan", file)
	assert.Nil(t, err)
	assert.True(t, u)
	assert.True(t, IsFileExisting(file))

	file = filepath.Join(tmp, "file1")
	u, err = Download("https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan", file)
	assert.Nil(t, err)
	assert.False(t, u)
	assert.True(t, IsFileExisting(file))

	err = os.WriteFile(file, []byte("some other data to change file"), 0644)
	assert.Nil(t, err)

	u, err = Download("https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan", file)
	assert.Nil(t, err)
	assert.True(t, u)
	assert.True(t, IsFileExisting(file))

	file = filepath.Join(tmp, "file2")
	u, err = Download("https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan_notexisting", file)
	assert.NotNil(t, err)
	assert.False(t, u)
	assert.False(t, IsFileExisting(file))
}

func TestUpdateDictionary(t *testing.T) {

}

// helper func
func prepareTempDir(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "toyskkserv-test")
	assert.Nil(t, err)

	return tmp
}
