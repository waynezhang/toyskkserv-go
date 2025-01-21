package dictionary

import (
	"bufio"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	getter "github.com/hashicorp/go-getter"
	cp "github.com/otiai10/copy"
	"github.com/waynezhang/tskks/internal/iconv"
	"github.com/waynezhang/tskks/internal/utils"
)

const dictEncodingPrefix = ";; -*- mode: fundamental; coding: "

func UpdateDictionaries(urls []string, dictDirectory string, cacheDirectory string) {
	slog.Info("Updating dictionaries", "urls", urls, "dictDirectory", dictDirectory)

	err := os.MkdirAll(dictDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, url := range urls {
		if err := DownloadDictionary(url, dictDirectory, cacheDirectory); err != nil {
			slog.Error("Failed to download file", "url", url)
		}
	}
}

func DownloadDictionary(url, dictDirectory string, cacheDirectory string) error {
	filename := dictName(url)
	path := filepath.Join(cacheDirectory, filename)
	dictPath := filepath.Join(dictDirectory, filename)

	slog.Info("Downloading file", "file", path)

	if err := getter.GetFile(path, url); err != nil {
		return err
	}

	srcChecksum := checksumOf(path)
	newSrcChecksum, _ := utils.FileChecksum(path)
	if srcChecksum != nil && *srcChecksum == newSrcChecksum && utils.IsFileExisting(dictPath) {
		slog.Info("Dict is up-to-date", "path", path)
		return nil
	}

	if err := updateUTF8Dictionary(path, dictPath); err != nil {
		slog.Error("Failed to update UTF-8 dictionary", "path", path, "err", err)
		return nil
	}

	setChecksumOf(path, newSrcChecksum)

	slog.Info("Dictionary is updated", "path", path)
	return nil
}

func updateUTF8Dictionary(src string, dst string) error {
	enc, err := encodingOfDict(src)
	if err != nil {
		return err
	}

	if enc == ENCODING_UTF8 {
		return cp.Copy(src, dst)
	}

	return convertEncoding(src, dst)
}

func convertEncoding(src string, dst string) error {
	slog.Info("Convert encoding", "src", src, "dst", dst)
	src_f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer src_f.Close()

	dst_f, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dst_f.Close()

	iv, err := iconv.Open(src, "euc-jis-2004", "utf-8")
	if err != nil {
		return err
	}
	defer iv.Close()

	w := bufio.NewWriter(dst_f)
	w.WriteString(";; -*- mode: fundamental; coding: utf-8 -*-\n")

	s := bufio.NewScanner(src_f)
	for s.Scan() {
		line := iv.Convert(s.Text())
		if strings.HasPrefix(line, dictEncodingPrefix) {
			continue
		}
		w.WriteString(line)
		w.WriteString("\n")
	}
	w.Flush()

	return nil
}

func encodingOfDict(path string) (encoding, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	if !s.Scan() {
		// fallback
		slog.Info("Dict file is empty", "file", path)
		return ENCODING_EUCJP, nil
	}

	// -*- mode: fundamental; coding: euc-jis-2004 -*-
	re := regexp.MustCompile(`coding:\s*([\w-]+)`)
	matches := re.FindStringSubmatch(s.Text())
	if len(matches) == 0 {
		slog.Info("No encoding definition found", "file", path)
		return ENCODING_EUCJP, nil
	}

	slog.Info("Encoding definition found", "encoding", matches[1], "file", path)
	if strings.HasPrefix(matches[1], "utf-8") {
		return ENCODING_UTF8, nil
	} else {
		return ENCODING_EUCJP, nil
	}
}

func checksumOf(path string) *string {
	sumfile := path + ".checksum"
	bytes, err := os.ReadFile(sumfile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		slog.Warn("Failed to read checksum file", "path", sumfile, "err", err)
		return nil
	}

	s := string(bytes)
	return &s
}

func setChecksumOf(path string, checksum string) error {
	sumfile := path + ".checksum"
	return os.WriteFile(sumfile, []byte(checksum), 0644)
}
