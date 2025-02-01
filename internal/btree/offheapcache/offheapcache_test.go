package offheapcache

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffheapCache(t *testing.T) {
	c := New(1024 * 1024)

	n := c.ReservedNode([]byte("a test"))
	assert.Equal(t, int32(0), n)

	fn := c.FindNode(0)
	k, _ := c.Bytes(0)
	assert.Equal(t, int(fn.keyLen), len("a test"))
	assert.Equal(t, []byte("a test"), k)

	cases := []string{
		"abc",
		"defefg",
		"xxxxxxxxx",
		"ed999",
		"23920234-^-^234",
	}

	adds := []int32{}
	for i, cse := range cases {
		_, add1 := c.NewNode(
			[]byte(cse),
			[]byte(cse+" - "+strconv.Itoa(i)),
		)
		adds = append(adds, add1)
	}

	for i, add := range adds {
		k, _ := c.Bytes(add)
		assert.Equal(t, cases[i], string(k))
	}

	c.Link(adds[1], adds[4])
	assert.Equal(t, c.FindNode(adds[1]).Next, adds[4])
}

func TestOffheapCacheFilling(t *testing.T) {
	keyOf := func(i int) []byte {
		return bytes.Repeat([]byte{'a' + byte(i%26)}, i*100)
	}
	valueOf := func(i int) []byte {
		return bytes.Repeat([]byte{'A' + byte(i%26)}, i*200)
	}

	c := New(1024 * 1024 * 100)
	offset := int32(reservedSpace + 16)
	assert.Equal(t, offset, c.offset)

	for i := 0; i < 100; i++ {
		k := keyOf(i)
		v := valueOf(i)
		_, _ = c.NewNode(k, v)
	}

	for i := 0; i < 100; i++ {
		magic := readInt32(c.buf[offset:])
		assert.Equal(t, magicNum, magic)
		offset += 4

		keyLen := readInt32(c.buf[offset:])
		assert.Equal(t, int32(i*100), keyLen)
		offset += 4

		valLen := readInt32(c.buf[offset:])
		assert.Equal(t, int32(i*200), valLen)
		offset += 4

		next := readInt32(c.buf[offset:])
		assert.Equal(t, int32(0), next)
		offset += 4

		expectedKeyBytes := bytes.Repeat([]byte{'a' + byte(i%26)}, int(keyLen))
		assert.Equal(t, expectedKeyBytes, c.buf[offset:offset+keyLen])
		offset += keyLen

		expectedValBytes := bytes.Repeat([]byte{'A' + byte(i%26)}, int(valLen))
		assert.Equal(t, expectedValBytes, c.buf[offset:offset+valLen])
		offset += valLen
	}

	assert.Equal(t, offset, c.offset)

	c.Clear()
	assert.Equal(t, int32(reservedSpace+16), c.offset)
}
