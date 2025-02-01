package dictionary

import (
	"io"
	"log/slog"
	"path/filepath"
	"sync"

	"github.com/waynezhang/toyskkserv/internal/dictionary/candidate"
	"github.com/waynezhang/toyskkserv/internal/files"
	"github.com/waynezhang/toyskkserv/internal/googleapi"
)

type DictManager struct {
	cm               *candidate.Manager
	directory        string
	fallbackToGoogle bool
}

var (
	instance *DictManager
	once     sync.Once
)

type Config struct {
	Dictionaires     []string
	Directory        string
	FallbackToGoogle bool
	UseDiskCache     bool
}

func NewDictManager(cfg Config) *DictManager {
	dm := &DictManager{
		cm:               candidate.New(cfg.UseDiskCache),
		directory:        cfg.Directory,
		fallbackToGoogle: cfg.FallbackToGoogle,
	}
	dm.reloadDicts(cfg.Dictionaires)

	return dm
}

func (dm *DictManager) HandleRequest(key string, w io.Writer) {
	slog.Info("Start finding candidates")
	defer slog.Info("Finished finding candidates")

	candidates := dm.cm.Find(key)

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
	dm.cm.IterateKey(key, func(c string) {
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
	dm.cm.Clear()

	dm.downloadDictionaries(urls)
	dm.loadFiles(files.DictionaryPaths(urls, dm.directory))
}

func (dm *DictManager) downloadDictionaries(urls []string) {
	slog.Info("Start downloading dictionaries")

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

	slog.Info("All dictionaries downloaded")
}

func (dm *DictManager) loadFiles(paths []string) {
	slog.Info("Start loading dictionaries")

	for _, path := range paths {
		loadFile(path, dm.cm)
	}

	slog.Info("All dictionaries loaded", "entries", dm.cm.Count())
}
