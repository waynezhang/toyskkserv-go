package dictionary

import (
	"log/slog"
	"strings"

	"github.com/google/btree"
)

type pair struct {
	key string
	val string
}

type candidatesManager struct {
	candidates *btree.BTreeG[*pair]
}

func newCandidatesManager() *candidatesManager {
	return &candidatesManager{
		candidates: btree.NewG(32, func(a, b *pair) bool { return a.key < b.key }),
	}
}

func (cm *candidatesManager) addCandidates(key string, candidates string) {
	p, ok := cm.candidates.Get(&pair{key: key})
	if ok {
		p.val = p.val + candidates[1:]
	} else {
		p = &pair{
			key: key,
			val: candidates,
		}
	}
	cm.candidates.ReplaceOrInsert(p)
}

func (cm *candidatesManager) findCandidates(key string) string {
	slog.Info("Find candidates", "key", "["+key+"]")
	pair := &pair{
		key: key,
	}
	if p, ok := cm.candidates.Get(pair); ok {
		return p.val
	}
	return ""
}

func (cm *candidatesManager) iterateCompletions(key string, ite func(c string)) {
	slog.Info("Find completions", "key", "["+key+"]")

	cm.candidates.AscendGreaterOrEqual(&pair{key: key}, func(p *pair) bool {
		if !strings.HasPrefix(p.key, key) {
			return false
		}

		ite(p.key)
		return true
	})
}

func (cm *candidatesManager) clear() {
	cm.candidates.Clear(false)
}
