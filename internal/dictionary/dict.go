package dictionary

import (
	"bufio"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type encoding string

const (
	ENCODING_UTF8  = "utf-8"
	ENCODING_EUCJP = "euc-jisx0213"
)

func loadDict(path string, cm *CandidatesManager) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	defer f.Close()

	if err != nil {
		slog.Error("Failed to open file", "file", path, "err", err)
		return
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		readLine(string(s.Bytes()), cm)
	}
}

func readLine(str string, cm *CandidatesManager) {
	if strings.HasPrefix(str, ";;") {
		return
	}

	keyEnd := strings.Index(str, " ")
	key := string(str[:keyEnd])
	candidates := string(str[keyEnd+1 : len(str)-1]) // start from /, trim the last slash
	cm.addCandidates(key, candidates)
}

func dictName(url string) string {
	return filepath.Base(url)
}
