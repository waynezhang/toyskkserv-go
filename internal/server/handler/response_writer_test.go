package handler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandidateRespnseWriter1(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	w := newCandidateResponseWriter(buf, "someword")
	w.Write([]byte("/cdd1/cdd2/"))
	w.wrap()

	assert.Equal(t, "1/cdd1/cdd2/\n", buf.String())
}

func TestRespnseWriter2(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	w := newCandidateResponseWriter(buf, "someword")
	w.wrap()

	assert.Equal(t, "4someword \n", buf.String())
}
