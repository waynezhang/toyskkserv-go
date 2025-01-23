package dictionary

import (
	"log/slog"
	"strings"

	"github.com/tidwall/btree"
)

type candidatesManager struct {
	candidates btree.Map[string, string]
}

func newCandidatesManager() *candidatesManager {
	return &candidatesManager{}
}

func (cm *candidatesManager) addCandidates(key string, candidates string) {
	old, _ := cm.candidates.Get(key)
	if len(old) == 0 {
		cm.candidates.Set(key, candidates)
	} else {
		cm.candidates.Set(key, old+candidates[1:])
	}
}

func (cm *candidatesManager) findCandidates(key string) string {
	slog.Info("Find candidates", "key", "["+key+"]")
	c, _ := cm.candidates.Get(key)
	return c
}

func (cm *candidatesManager) findCompletions(key string) string {
	slog.Info("Find completions", "key", "["+key+"]")

	cdd := ""
	cm.candidates.Ascend(key, func(k string, v string) bool {
		if !strings.HasPrefix(k, key) {
			return false
		}
		cdd += k + "/"
		return true
	})
	if len(cdd) > 0 {
		cdd = "/" + cdd
	}
	return cdd
}

func (cm *candidatesManager) clear() {
	cm.candidates.Clear()
}
