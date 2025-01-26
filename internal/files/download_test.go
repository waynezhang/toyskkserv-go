package files

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFunction struct {
	mock.Mock
	dontCreateFile bool
}

func (m *MockFunction) download(url string, dst string) error {
	if !m.dontCreateFile {
		os.WriteFile(dst, []byte("test data"), 0644)
	}
	args := m.Called(url, dst)
	if args.Get(0) == nil {
		return nil
	} else {
		return args.Error(0)
	}
}

func (m *MockFunction) notify() {
	m.Called()
}

func TestDownload(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	url := "https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.edict"
	path := filepath.Join(tmp, "dest_file")

	mockObj := new(MockFunction)

	// new downloaded case
	mockObj.On("download", url, path).Return(nil)
	updated, err := download(url, path, mockObj.download)
	assert.True(t, updated)
	assert.Nil(t, err)

	// updated case
	assert.Nil(t, os.RemoveAll(path))
	assert.Nil(t, os.WriteFile(path, []byte("original data"), 0644))

	mockObj.On("download", url, path).Return(nil)
	updated, err = download(url, path, mockObj.download)
	assert.True(t, updated)
	assert.Nil(t, err)

	// not updated case
	assert.Nil(t, os.RemoveAll(path))
	assert.Nil(t, os.WriteFile(path, []byte("test data"), 0644))

	mockObj.On("download", url, path).Return(nil)
	updated, err = download(url, path, mockObj.download)
	assert.False(t, updated)
	assert.Nil(t, err)

	// download succeed but no file created
	assert.Nil(t, os.RemoveAll(path))

	mockObj.dontCreateFile = true
	mockObj.On("download", url, path).Return(errors.New("some error"))
	updated, err = download(url, path, mockObj.download)
	assert.False(t, updated)
	assert.NotNil(t, err)

	// download failed
	assert.Nil(t, os.RemoveAll(path))
	assert.Nil(t, os.WriteFile(path, []byte("test data"), 0644))

	mockObj.On("download", mock.Anything, mock.Anything).Unset()
	mockObj.On("download", url, path).Return(errors.New("some error"))
	updated, err = download(url, path, mockObj.download)
	assert.False(t, updated)
	assert.NotNil(t, err)
}

func TestUpdateDictionary(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	mockObj := new(MockFunction)
	mockObj.On("notify").Return()

	urls := []string{
		"https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.edict",
		"https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan",
		filepath.Join(tmp, "some_local_file"),
	}
	paths := []string{
		filepath.Join(tmp, "SKK-JISYO.edict"),
		filepath.Join(tmp, "SKK-JISYO.china_taiwan"),
		filepath.Join(tmp, "some_local_file"),
	}
	assert.Nil(t, os.WriteFile(paths[0], []byte("original data"), 0644))
	assert.Nil(t, os.WriteFile(paths[1], []byte("original data"), 0644))
	assert.Nil(t, os.WriteFile(paths[2], []byte("original data"), 0644))

	mockObj.On("download", urls[0], paths[0]).Return(nil)
	mockObj.On("download", urls[1], paths[1]).Return(nil)
	mockObj.On("download", urls[2], paths[2]).Panic("not supposed to be called")

	updateDictionaries(urls, tmp, mockObj.download, mockObj.notify)

	mockObj.AssertCalled(t, "download", urls[0], paths[0])
	mockObj.AssertCalled(t, "download", urls[1], paths[1])
	mockObj.AssertCalled(t, "notify")

	// no update case

	mockObj2 := new(MockFunction)
	mockObj2.On("download", mock.Anything, mock.Anything).Return(nil)
	mockObj2.On("notify").Return(nil)
	assert.Nil(t, os.WriteFile(paths[0], []byte("test data"), 0644))
	assert.Nil(t, os.WriteFile(paths[1], []byte("test data"), 0644))
	assert.Nil(t, os.WriteFile(paths[2], []byte("test data"), 0644))

	updateDictionaries(urls, tmp, mockObj2.download, mockObj2.notify)

	mockObj2.AssertCalled(t, "download", urls[0], paths[0])
	mockObj2.AssertCalled(t, "download", urls[1], paths[1])
	mockObj2.AssertNotCalled(t, "notify")
}

func TestRealDownload(t *testing.T) {
	tmp := prepareTempDir(t)
	defer os.RemoveAll(tmp)

	urls := []string{
		"https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.edict",
		"https://github.com/skk-dev/dict/raw/refs/heads/master/SKK-JISYO.china_taiwan",
		filepath.Join(tmp, "some_local_file"),
	}

	mockObj := new(MockFunction)
	mockObj.On("notify").Return()

	updateDictionaries(urls, tmp, httpDownload, mockObj.notify)

	assert.True(t, IsFileExisting(filepath.Join(tmp, "SKK-JISYO.edict")))
	assert.True(t, IsFileExisting(filepath.Join(tmp, "SKK-JISYO.china_taiwan")))
	mockObj.AssertCalled(t, "notify")
}

// helper func
func prepareTempDir(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "toyskkserv-test")
	assert.Nil(t, err)

	return tmp
}
