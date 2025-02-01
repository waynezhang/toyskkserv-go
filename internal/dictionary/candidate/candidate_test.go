package candidate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/toyskkserv/internal/btree"
)

func TestCandidates(t *testing.T) {
	cases := []bool{true, false}
	for _, c := range cases {
		m := New(c)
		if c {
			assert.IsType(t, &btree.OffheapBTree{}, m.tree)
		} else {
			assert.IsType(t, &btree.InMemBTree{}, m.tree)
		}

		assert.Equal(t, "", m.Find("abc"))

		m.Add("abc", "/val1/")
		assert.Equal(t, "/val1/", m.Find("abc"))

		m.Add("abc", "/val2/val3/val4/")
		assert.Equal(t, "/val1/val2/val3/val4/", m.Find("abc"))

		m.Add("ABC", "/val3/")
		m.Add("ABC", "/val4/")
		assert.Equal(t, "/val1/val2/val3/val4/", m.Find("abc"))
		assert.Equal(t, "/val3/val4/", m.Find("ABC"))

		m.Clear()
		assert.Equal(t, "", m.Find("abc"))
		assert.Equal(t, "", m.Find("ABC"))
	}
}

func TestCompletions(t *testing.T) {
	cases := []bool{true, false}
	for _, c := range cases {
		m := New(c)

		m.Add("abc", "/val1/")
		m.Add("abc", "/val2/")
		m.Add("ABC", "/val3/")
		m.Add("ABC", "/val4/")
		m.Add("abd", "/val3/")
		m.Add("abd", "/val4/")

		cases := [][]string{
			{"a", "/abc/abd/"},
			{"ab", "/abc/abd/"},
			{"abc", "/abc/"},
			{"A", "/ABC/"},
			{"def", ""},
		}

		for _, c := range cases {
			key := c[0]
			comps := ""
			m.IterateKey(key, func(c string) {
				comps += c + "/"
			})
			if len(comps) > 0 {
				comps = "/" + comps
			}
			assert.Equal(t, c[1], comps)
		}
	}
}
