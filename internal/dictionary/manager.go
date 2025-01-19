package dictionary

import (
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	"github.com/waynezhang/tskks/internal/config"
	"github.com/waynezhang/tskks/internal/utils"
)

type DictManager struct {
	Directory string
	cm        *CandidatesManager
}

var (
	instance *DictManager
	once     sync.Once
)

func Shared() *DictManager {
	once.Do(func() {
		cfg := config.Shared()
		instance = &DictManager{
			cm:        newCandidatesManager(),
			Directory: cfg.DictionaryDirectory,
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
	} else {
		return "4/" + key + " "
	}
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
