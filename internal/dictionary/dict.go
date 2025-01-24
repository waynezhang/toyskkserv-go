package dictionary

import (
	"bufio"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/waynezhang/tskks/internal/iconv"
)

func loadFile(path string, cm *candidatesManager) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		slog.Error("Failed to open file", "file", path, "err", err)
		return
	}
	defer f.Close()

	enc := iconv.ENCODING_UNDECIDED
	if strings.HasSuffix(path, ".utf8") {
		enc = iconv.ENCODING_UTF8
	}

	s := bufio.NewScanner(f)
	idx := 0
	for s.Scan() {
		idx++
		line := string(s.Bytes())
		if strings.HasPrefix(line, ";;") {
			if enc == iconv.ENCODING_UNDECIDED {
				enc = parseEncoding(line)
			}
			continue
		}

		if enc == iconv.ENCODING_UNDECIDED {
			// default
			enc = iconv.ENCODING_EUCJP
		}
		if enc == iconv.ENCODING_UTF8 {
			parseLine(line, cm)
		} else {
			s, err := iconv.EUCJPConverter.ConvertLine(line)
			if err != nil {
				slog.Error("Failed to covnert encoding", "file", path, "no", idx, "text", line)
				continue
			}
			parseLine(s, cm)
		}
	}
}

func parseEncoding(str string) string {
	// -*- coding: euc-jis-2004 -*-
	re := regexp.MustCompile(`coding:\s*([\w-]+)`)
	matches := re.FindStringSubmatch(str)
	if len(matches) == 0 {
		return iconv.ENCODING_UNDECIDED
	}

	if strings.HasPrefix(matches[1], "utf-8") {
		return iconv.ENCODING_UTF8
	} else {
		// by default
		return iconv.ENCODING_EUCJP
	}
}

func parseLine(str string, cm *candidatesManager) {
	if strings.HasPrefix(str, ";;") {
		return
	}

	keyEnd := strings.Index(str, " ")
	key := string(str[:keyEnd])
	candidates := string(str[keyEnd+1:]) // /val1/
	cm.addCandidates(key, candidates)
}
