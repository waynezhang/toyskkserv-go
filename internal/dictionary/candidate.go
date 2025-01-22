package dictionary

import (
	"log/slog"
	"strings"

	"github.com/tidwall/btree"
)

type CandidatesManager struct {
	candidates btree.Map[string, string]
}

func newCandidatesManager() *CandidatesManager {
	return &CandidatesManager{}
}

func (cm *CandidatesManager) addCandidates(key string, candidates string) {
	old, _ := cm.candidates.Get(key)
	cm.candidates.Set(key, old+candidates)
}

func (cm *CandidatesManager) findCandidates(key string) string {
	slog.Info("Find candidates", "key", "["+key+"]")
	c, _ := cm.candidates.Get(key)
	return c
}

func (cm *CandidatesManager) findCompletions(key string) string {
	slog.Info("Find completions", "key", "["+key+"]")

	cdd := ""
	cm.candidates.Ascend(key, func(k string, v string) bool {
		if !strings.HasPrefix(k, key) {
			return false
		}
		cdd += "/" + k
		return true
	})
	return cdd
}

func (cm *CandidatesManager) clear() {
	cm.candidates.Clear()
}
