package dictionary

import (
	"io"
	"log/slog"
	"path/filepath"
	"sync"

	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/files"
	"github.com/waynezhang/toyskkserv/internal/googleapi"
)

type DictManager struct {
	cm               *candidatesManager
	directory        string
	fallbackToGoogle bool
}

var (
	instance *DictManager
	once     sync.Once
)

func Shared() *DictManager {
	once.Do(func() {
		cfg := config.Shared()

		dm := NewDictManager(cfg.DictionaryDirectory, cfg.FallbackToGoogle)
		dm.reloadDicts(cfg.Dictionaries)

		instance = dm
	})
	return instance
}

func NewDictManager(directory string, fallbackToGoogle bool) *DictManager {
	return &DictManager{
		cm:               newCandidatesManager(),
		directory:        directory,
		fallbackToGoogle: fallbackToGoogle,
	}
}

func (dm *DictManager) HandleRequest(key string, w io.Writer) {
	slog.Info("Start finding candidates")
	defer slog.Info("Finished finding candidates")

	candidates := dm.cm.findCandidates(key)

	if len(candidates) > 0 {
		w.Write([]byte(candidates))
		return
	}

	if dm.fallbackToGoogle {
		candidates = googleapi.TransliterateRequest(key)
		if len(candidates) > 0 {
			w.Write([]byte(candidates))
			w.Write([]byte{'/'})
			return
		}
	}
}

func (dm *DictManager) HandleCompletion(key string, w io.Writer) {
	slog.Info("Start finding completions")
	defer slog.Info("Finished finding completions")

	found := false
	dm.cm.iterateCompletions(key, func(c string) {
		found = true
		w.Write([]byte{'/'})
		w.Write([]byte([]byte(c)))
	})
	if found {
		w.Write([]byte{'/'})
	}
}

func (dm *DictManager) DictionariesDidChange(urls []string) {
	slog.Info("Dictionaries did change")
	dm.reloadDicts(urls)
}

func (dm *DictManager) reloadDicts(urls []string) {
	dm.cm.clear()

	dm.downloadDictionaries(urls)
	dm.loadFiles(files.DictionaryPaths(urls, dm.directory))
}

func (dm *DictManager) downloadDictionaries(urls []string) {
	slog.Info("Start loading dictionaries")

	for _, url := range urls {
		if files.IsLocalURL(url) {
			slog.Info("Skip local file", "url", url)
			continue
		}

		path := filepath.Join(dm.directory, files.DictName(url))
		if files.IsFileExisting(path) {
			continue
		}
		slog.Warn("Dictionary not found", "path", path)
		files.Download(url, path)
	}

	slog.Info("All dictionaries loaded")
}

func (dm *DictManager) loadFiles(paths []string) {
	for _, path := range paths {
		loadFile(path, dm.cm)
	}
}
