package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
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
