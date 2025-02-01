package btree

import (
	"strings"

	"github.com/google/btree"
)

type inMemTreePair struct {
	key string
	val string
}

type InMemBTree struct {
	t *btree.BTreeG[inMemTreePair]
}

func NewInMemTree() *InMemBTree {
	return &InMemBTree{
		t: btree.NewG(32, func(a, b inMemTreePair) bool {
			return a.key < b.key
		}),
	}
}

func (t *InMemBTree) Get(key string) (string, bool) {
	if p, ok := t.t.Get(inMemTreePair{key: key}); ok {
		return p.val, true
	}

	return "", false
}

func (t *InMemBTree) Append(key, existingVal, val string) {
	t.t.ReplaceOrInsert(inMemTreePair{
		key: key,
		val: existingVal + val,
	})
}

func (t *InMemBTree) IterateKey(prefix string, fn func(key string)) {
	t.t.AscendGreaterOrEqual(inMemTreePair{key: prefix}, func(p inMemTreePair) bool {
		if !strings.HasPrefix(p.key, prefix) {
			return false
		}

		fn(p.key)
		return true
	})
}

func (t *InMemBTree) Count() int {
	return t.t.Len()
}

func (t *InMemBTree) Clear() {
	t.t.Clear(true)
}
