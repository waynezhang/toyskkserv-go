package dictionary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandidates(t *testing.T) {
	cm := newCandidatesManager()

	assert.Equal(t, "", cm.findCandidates("abc"))

	cm.addCandidates("abc", "/val1/")
	assert.Equal(t, "/val1/", cm.findCandidates("abc"))

	cm.addCandidates("abc", "/val2/val3/val4/")
	assert.Equal(t, "/val1/val2/val3/val4/", cm.findCandidates("abc"))

	cm.addCandidates("ABC", "/val3/")
	cm.addCandidates("ABC", "/val4/")
	assert.Equal(t, "/val1/val2/val3/val4/", cm.findCandidates("abc"))
	assert.Equal(t, "/val3/val4/", cm.findCandidates("ABC"))

	cm.clear()
	assert.Equal(t, "", cm.findCandidates("abc"))
	assert.Equal(t, "", cm.findCandidates("ABC"))
}

func TestCompletions(t *testing.T) {
	cm := newCandidatesManager()

	cm.addCandidates("abc", "/val1/")
	cm.addCandidates("abc", "/val2/")
	cm.addCandidates("ABC", "/val3/")
	cm.addCandidates("ABC", "/val4/")
	cm.addCandidates("abd", "/val3/")
	cm.addCandidates("abd", "/val4/")

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
		cm.iterateCompletions(key, func(c string) {
			comps += c + "/"
		})
		if len(comps) > 0 {
			comps = "/" + comps
		}
		assert.Equal(t, c[1], comps)
	}
}
