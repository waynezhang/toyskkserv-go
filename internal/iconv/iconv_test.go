package iconv

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIconv(t *testing.T) {
	iv, err := open(ENCODING_EUCJP, ENCODING_UTF8)

	assert.Nil(t, err)
	defer iv.close()

	f, err := os.Open("../../testdata/jisyo.euc-jp")
	r := bufio.NewScanner(f)
	read := false
	for r.Scan() {
		read = true
		line := r.Text()
		_, err := iv.ConvertLine(line)
		assert.Nil(t, err)
	}

	assert.True(t, read)
}

func TestConverter(t *testing.T) {
	cases := [][]string{
		{"1234", "1234"},
	}
	for _, c := range cases {
		s, err := EUCJPConverter.ConvertLine(c[0])
		assert.Nil(t, err)
		assert.Equal(t, c[0], s)
	}
}
