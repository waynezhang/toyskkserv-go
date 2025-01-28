package dictionary

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/waynezhang/eucjis2004decode/eucjis2004"
	"golang.org/x/text/transform"
)

const (
	ENCODING_UNDECIDED = "undecided"
	ENCODING_UTF8      = "utf-8"
	ENCODING_EUCJP     = "euc-jp"
)

func loadFile(path string, cm *candidatesManager) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		slog.Error("Failed to open file", "file", path, "err", err)
		return
	}
	defer f.Close()

	enc := ENCODING_UNDECIDED

	if strings.HasSuffix(path, ".utf8") {
		enc = ENCODING_UTF8
	} else {
		enc, err = detectFileEncoding(f)
		if err != nil {
			slog.Error("Failed to detect encoding", "file", path, "err", err)
			return
		}
	}

	var s *bufio.Scanner
	if enc == ENCODING_UTF8 {
		s = bufio.NewScanner(f)
	} else {
		s = bufio.NewScanner(transform.NewReader(f, eucjis2004.EUCJIS2004Decoder{}))
	}

	for s.Scan() {
		parseLine(s.Bytes(), cm)
	}
}

func detectFileEncoding(f *os.File) (string, error) {
	currPos, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		}
		return ENCODING_EUCJP, err
	}

	defer f.Seek(currPos, 0)

	// ;; -*- mode: fundamental; coding: euc-jp -*-
	b := make([]byte, 1024)
	_, err = f.Read(b)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		}
		return ENCODING_EUCJP, err
	}

	coding := []byte("coding: ")
	idx := bytes.Index(b, coding)
	if idx == -1 {
		return ENCODING_EUCJP, err
	}

	if bytes.HasPrefix(b[idx+len(coding):], []byte("utf-8")) {
		return ENCODING_UTF8, nil
	} else {
		// by default
		return ENCODING_EUCJP, nil
	}
}

func parseLine(bs []byte, cm *candidatesManager) {
	if bytes.HasPrefix(bs, []byte{';', ';'}) {
		return
	}

	keyEnd := bytes.IndexByte(bs, ' ')
	key := string(bs[:keyEnd])
	candidates := string(bs[keyEnd+1:]) // /val1/
	cm.addCandidates(key, candidates)
}
