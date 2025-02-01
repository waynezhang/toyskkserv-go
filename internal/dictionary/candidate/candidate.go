package candidate

import (
	"log/slog"
	"sync"

	"github.com/waynezhang/toyskkserv/internal/btree"
)

type Manager struct {
	tree     btree.BTree
	mu       sync.Mutex
	updating bool
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
	if !m.updating {
		return
	}

	if p, ok := m.tree.Get(key); ok {
		m.tree.Append(key, p, candidates[1:])
	} else {
		m.tree.Append(key, "", candidates)
	}

}

func (m *Manager) Find(key string) string {
	if m.updating {
		return ""
	}

	slog.Info("Find candidates", "key", "["+key+"]")

	if p, ok := m.tree.Get(key); ok {
		return p
	}
	return ""
}

func (m *Manager) IterateKey(key string, ite func(c string)) {
	if m.updating {
		return
	}

	slog.Info("Find completions", "key", "["+key+"]")
	m.tree.IterateKey(key, ite)
}

func (m *Manager) Count() int {
	if m.updating {
		return 0
	}

	return m.tree.Count()
}

func (m *Manager) Transaction(fn func(m *Manager)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.updating = true
	fn(m)
	m.updating = false
}

func (m *Manager) Clear() {
	if !m.updating {
		return
	}

	m.tree.Clear()
}
