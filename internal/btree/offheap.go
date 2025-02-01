package btree

import (
	"bytes"
	"strings"
	"sync"

	"github.com/google/btree"
	"github.com/waynezhang/toyskkserv/internal/btree/offheapcache"
)

type OffheapBTree struct {
	t     *btree.BTreeG[int32]
	cache *offheapcache.Cache
	mu    sync.Mutex
}

func NewOffheapBtree() *OffheapBTree {
	cache := offheapcache.New(0)
	if cache == nil {
		return nil
	}

	return &OffheapBTree{
		t: btree.NewG(32, func(a, b int32) bool {
			akey, _ := cache.Bytes(a)
			bkey, _ := cache.Bytes(b)
			return bytes.Compare(akey, bkey) == -1
		}),
		cache: cache,
	}
}

func (t *OffheapBTree) Get(key string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	resvAddr := t.cache.ReservedNode([]byte(key))
	if p, ok := t.t.Get(resvAddr); ok {
		addr := p
		ret := ""
		for {
			_, val := t.cache.Bytes(addr)
			ret = ret + string(val)

			n := t.cache.FindNode(addr)
			if n.Next != 0 {
				addr = n.Next
				continue
			}

			return ret, true
		}
	}
	return "", false
}

func (t *OffheapBTree) Append(key, existingVal, val string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if existingVal == "" {
		_, addr := t.cache.NewNode([]byte(key), []byte(val))
		t.t.ReplaceOrInsert(addr)
		return
	}

	resvAddr := t.cache.ReservedNode([]byte(key))
	p, ok := t.t.Get(resvAddr)
	if !ok {
		panic(key)
	}

	_, newAddr := t.cache.NewNode([]byte(key), []byte(val))
	addr := p
	for {
		n := t.cache.FindNode(addr)
		if n.Next != 0 {
			addr = n.Next
			continue
		}

		t.cache.Link(addr, newAddr)
		break
	}
}

func (t *OffheapBTree) IterateKey(prefix string, fn func(key string)) {
	t.mu.Lock()
	defer t.mu.Unlock()

	resvAddr := t.cache.ReservedNode([]byte(prefix))
	t.t.AscendGreaterOrEqual(resvAddr, func(p int32) bool {
		k, _ := t.cache.Bytes(p)
		if !strings.HasPrefix(string(k), prefix) {
			return false
		}

		fn(string(k))
		return true
	})
}

func (t *OffheapBTree) Count() int {
	return t.t.Len()
}

func (t *OffheapBTree) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.t.Clear(true)
	t.cache.Clear()
}
