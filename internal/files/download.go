package files

import (
	"log/slog"
	"os"

	getter "github.com/hashicorp/go-getter"
	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/tcp"
)

func UpdateDictionaries(urls []string, dictDirectory string) {
	slog.Info("Updating dictionaries", "urls", urls, "dictDirectory", dictDirectory)

	err := os.MkdirAll(dictDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}

	updated := false
	for _, url := range urls {
		u, err := Download(url, dictDirectory)
		if err != nil {
			slog.Error("Failed to download file", "url", url)
			continue
		}
		updated = updated || u
	}

	if updated {
		tcp.SendReloadCommand(config.Shared().ListenAddr)
	}
}

func Download(url, tofile string) (updated bool, err error) {
	oldChecksum, err := FileChecksum(tofile)
	if err != nil && IsFileExisting(tofile) {
		slog.Error("Failed to get checksum from existing file", "path", tofile, "err", err)
	}

	slog.Info("Downloading file", "file", tofile)
	if err := getter.GetFile(tofile, url); err != nil {
		os.Remove(tofile)
		return false, err
	}

	newChecksum, err := FileChecksum(tofile)
	if err != nil {
		slog.Error("Failed to get checksum from new downloaded file", "path", tofile, "err", err)
	}

	if oldChecksum == newChecksum {
		slog.Info("Dict is up-to-date", "path", tofile)
		return false, nil
	}

	slog.Info("File is updated", "path", tofile)
	return true, nil
}
