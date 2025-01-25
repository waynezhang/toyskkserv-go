package dictionary

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/waynezhang/eucjis2004decode/decode"
)

const (
	ENCODING_UNDECIDED = "undecided"
	ENCODING_UTF8      = "utf-8"
	ENCODING_EUCJP     = "euc-jp"
)

// -*- coding: euc-jis-2004 -*-
var codingRegex = regexp.MustCompile(`coding:\s*([\w-]+)`)

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
	}

	buf := bytes.NewBuffer(nil)

	s := bufio.NewScanner(f)
	idx := 0
	for s.Scan() {
		idx++
		buf.Reset()

		bs := s.Bytes()
		if bytes.HasPrefix(bs, []byte{';', ';'}) {
			if enc == ENCODING_UNDECIDED {
				enc = parseEncoding(bs)
			}
			continue
		}

		if enc == ENCODING_UNDECIDED {
			// default
			enc = ENCODING_EUCJP
		}
		if enc == ENCODING_UTF8 {
			parseLine(bs, cm)
		} else {
			err := decode.Convert(s.Bytes(), buf)
			if err != nil {
				slog.Error("Failed to covnert encoding", "file", path, "no", idx, "text", string(bs))
				continue
			}
			parseLine(buf.Bytes(), cm)
		}
	}

	buf.Reset()
}

func parseEncoding(bs []byte) string {
	matches := codingRegex.FindSubmatch(bs)
	if len(matches) == 0 {
		return ENCODING_UNDECIDED
	}

	if bytes.Compare(matches[1], []byte{'u', 't', 'f', '-', '8'}) == 0 {
		return ENCODING_UTF8
	} else {
		// by default
		return ENCODING_EUCJP
	}
}

func parseLine(bs []byte, cm *candidatesManager) {
	keyEnd := bytes.IndexByte(bs, ' ')
	key := string(bs[:keyEnd])
	candidates := string(bs[keyEnd+1:]) // /val1/
	cm.addCandidates(key, candidates)
}
