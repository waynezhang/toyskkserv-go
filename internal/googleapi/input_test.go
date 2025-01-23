package googleapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResponse(t *testing.T) {
	cases := [][]string{
		{"", "", ""},
		{"", "xxx", ""},
		{"", "[]", ""},
		{"", "[[]]", ""},
		{"", "[[\"a\"]]", ""},
		{"akey", "[[\"anotherkey\"]]", ""},
		{"akey", "[[\"anotherkey\", [\"v1\", \"v2\"]]]", ""},
		{"akey", "[[\"akey\", [\"akey\", \"v1\", \"v2\"]]]", "/v1/v2"},
	}

	for _, c := range cases {
		key := c[0]
		body := c[1]
		expected := c[2]

		ret := paraseResponse([]byte(body), key)
		assert.Equal(t, expected, ret)
	}
}
