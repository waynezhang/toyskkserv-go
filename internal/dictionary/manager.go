package dictionary

import (
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	"github.com/waynezhang/tskks/internal/config"
	"github.com/waynezhang/tskks/internal/googleapi"
	"github.com/waynezhang/tskks/internal/utils"
)

type DictManager struct {
	directory        string
	cm               *CandidatesManager
	fallbackToGoogle bool
}

var (
	instance *DictManager
	once     sync.Once
)

func Shared() *DictManager {
	once.Do(func() {
		cfg := config.Shared()
		instance = &DictManager{
			cm:               newCandidatesManager(),
			directory:        cfg.DictionaryDirectory,
			fallbackToGoogle: cfg.FallbackToGoogle,
		}
		instance.loadAll(cfg)
	})
	return instance
}

func (dm *DictManager) HandleRequest(req string) string {
	slog.Info("Start finding candidates")
	defer slog.Info("Finished finding candidates")

	key := strings.TrimSuffix(strings.TrimPrefix(req, "1"), " ")
	candidates := dm.cm.findCandidates(key)

	if len(candidates) > 0 {
		return "1" + candidates + "/"
	}

	if dm.fallbackToGoogle {
		candidates = googleapi.TransliterateRequest(key)
		if len(candidates) > 0 {
			return "1" + candidates + "/"
		}
	}

	return "4/" + key + " "
}

func (dm *DictManager) DictionariesDidChange() {
	dm.reloadDicts()
}

func (dm *DictManager) loadAll(cfg *config.Config) {
	slog.Info("Start loading dictionaries")

	for _, url := range cfg.Dictionaries {
		dictPath := filepath.Join(cfg.DictionaryDirectory, dictName(url))
		if !utils.IsFileExisting(dictPath) {
			slog.Warn("Dictionary not found", "path", dictPath)
			DownloadDictionary(url, cfg.DictionaryDirectory, cfg.CacheDirectory)
		}
		loadDict(dictPath, dm.cm)
	}

	slog.Info("All dictionaries loaded")
}

func (dm *DictManager) reloadDicts() {
	dm.cm.candidates.Clear()
	dm.loadAll(config.Shared())
}
