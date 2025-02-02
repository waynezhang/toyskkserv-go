package btree

import (
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffheapBTree(t *testing.T) {
	tree := NewOffheapBtree()

	for i := 0; i < 10; i++ {
		k := strconv.Itoa(i % 2)
		v, _ := tree.Get(k)
		tree.Append(k, v, strconv.Itoa(i))
	}

	v0, _ := tree.Get("0")
	v1, _ := tree.Get("1")
	assert.Equal(t, "02468", v0)
	assert.Equal(t, "13579", v1)
}

func TestConcurrentAccess(t *testing.T) {
	tree := NewOffheapBtree()

	var wg sync.WaitGroup

	keys := make([]string, 100)
	vals := make([]string, 100)

	for i := 0; i < 100; i++ {
		keys[i] = strings.Repeat(string('a'+byte(i%10)), i*100)
		vals[i] = strings.Repeat(string('a'+byte((i+1)%10)), i*200)
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			tree.Append(keys[i], "", vals[i])
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 100, tree.Count())

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			v, ok := tree.Get(keys[i])
			assert.True(t, ok)
			assert.Equal(t, vals[i], v)
			wg.Done()
		}()
	}
	wg.Wait()

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			v, ok := tree.Get(keys[i])
			assert.True(t, ok)
			assert.Equal(t, vals[i], v)
			wg.Done()
		}()
	}
	wg.Wait()

	wg.Add(9)
	for i := 1; i < 10; i++ {
		go func() {
			count := 0
			tree.IterateKey(keys[i], func(key string) {
				count++
			})
			assert.Equal(t, 10, count)
			wg.Done()
		}()
	}
	wg.Wait()
}
