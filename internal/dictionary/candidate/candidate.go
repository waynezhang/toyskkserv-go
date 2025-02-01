package candidate

import (
	"log/slog"

	"github.com/waynezhang/toyskkserv/internal/btree"
)

type Manager struct {
	tree btree.BTree
}

func New(onDisk bool) *Manager {
	slog.Info("Creating candidates manager", "OnDisk", onDisk)

	if onDisk {
		t := btree.NewOffheapBtree()
		if t != nil {
			return &Manager{
				tree: t,
			}
		}
		slog.Warn("Fallback to memory cache")
	}
	return &Manager{
		tree: btree.NewInMemTree(),
	}
}

func (m *Manager) Add(key string, candidates string) {
	if p, ok := m.tree.Get(key); ok {
		m.tree.Append(key, p, candidates[1:])
	} else {
		m.tree.Append(key, "", candidates)
	}
}

func (m *Manager) Find(key string) string {
	slog.Info("Find candidates", "key", "["+key+"]")

	if p, ok := m.tree.Get(key); ok {
		return p
	}
	return ""
}

func (m *Manager) IterateKey(key string, ite func(c string)) {
	slog.Info("Find completions", "key", "["+key+"]")
	m.tree.IterateKey(key, ite)
}

func (m *Manager) Count() int {
	return m.tree.Count()
}

func (m *Manager) Clear() {
	m.tree.Clear()
}
