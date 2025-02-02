package responsewriter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandidateRespnseWriter1(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	w1 := New(buf, "someword")
	_, err := w1.Write([]byte("/cdd1/cdd2/"))
	assert.Nil(t, err)
	w1.Wrap()

	assert.True(t, w1.disposed)
	assert.Equal(t, "1/cdd1/cdd2/\n", buf.String())

	buf.Reset()

	w2 := New(buf, "anotherword")
	assert.False(t, w2.disposed)
	assert.Equal(t, w1, w2)

	_, err = w2.Write([]byte("/cdd3/cdd4/"))
	assert.Nil(t, err)
	w2.Wrap()

	assert.Equal(t, "1/cdd3/cdd4/\n", buf.String())
}

func TestRespnseWriter2(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	w := New(buf, "someword")
	w.Wrap()

	assert.Equal(t, "4someword \n", buf.String())
}
