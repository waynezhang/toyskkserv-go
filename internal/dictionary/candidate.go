package dictionary

import (
	"log/slog"

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

func (cm *CandidatesManager) clear() {
	cm.candidates.Clear()
}
