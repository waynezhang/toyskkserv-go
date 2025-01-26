package files

import (
	"log/slog"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/tcp"
)

type downloadFn func(url string, dst string) error
type serverNotifyFn func()

func httpDownload(url string, dst string) error {
	return getter.GetFile(dst, url)
}

func tcpServerNotifyFn() {
	tcp.SendReloadCommand(config.Shared().ListenAddr)
}

func UpdateDictionaries(urls []string, dictDirectory string) {
	updateDictionaries(urls, dictDirectory, httpDownload, tcpServerNotifyFn)
}

func updateDictionaries(urls []string, dictDirectory string, dnFn downloadFn, notifyFn serverNotifyFn) {
	slog.Info("Updating dictionaries", "urls", urls, "dictDirectory", dictDirectory)

	err := os.MkdirAll(dictDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}

	updated := false
	for _, url := range urls {
		if IsLocalURL(url) {
			slog.Info("Skip local file", "url", url)
			continue
		}

		path := filepath.Join(dictDirectory, filepath.Base(url))
		u, err := download(url, path, dnFn)
		if err != nil {
			slog.Error("Failed to download file", "url", url)
			continue
		}

		updated = updated || u
	}

	if updated {
		notifyFn()
	}
}

func Download(url string, tofile string) (updated bool, err error) {
	return download(url, tofile, httpDownload)
}

func download(url string, tofile string, dnFn downloadFn) (updated bool, err error) {
	oldChecksum, err := FileChecksum(tofile)
	if err != nil && IsFileExisting(tofile) {
		slog.Error("Failed to get checksum from existing file", "path", tofile, "err", err)
	}

	slog.Info("Downloading file", "file", tofile)
	if err := dnFn(url, tofile); err != nil {
		os.Remove(tofile)
		return false, err
	}

	newChecksum, err := FileChecksum(tofile)
	if err != nil {
		slog.Error("Failed to get checksum from new downloaded file", "path", tofile, "err", err)
		os.Remove(tofile)
		return false, err
	}

	if oldChecksum == newChecksum {
		slog.Info("Dict is up-to-date", "path", tofile)
		return false, nil
	}

	slog.Info("File is updated", "path", tofile)
	return true, nil
}
