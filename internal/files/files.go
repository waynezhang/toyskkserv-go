package files

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func FileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func IsFileExisting(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DictionaryPaths(urls []string, directory string) []string {
	paths := []string{}
	for _, u := range urls {
		if IsLocalURL(u) {
			e, _ := homedir.Expand(u)
			paths = append(paths, e)
			continue
		}
		p := filepath.Join(directory, DictName(u))
		paths = append(paths, p)
	}

	return paths
}

func DictName(url string) string {
	return filepath.Base(url)
}

func IsLocalURL(url string) bool {
	return !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://")
}
