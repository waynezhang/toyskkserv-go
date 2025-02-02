package offheapcache

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/edsrzf/mmap-go"
)

const (
	magicNum      = int32(0x09FFFFFF)
	reservedSpace = 1024 * 10
	tmpFileSize   = 1024 * 1024 * 100
	nodeSize      = 16
)

type node struct {
	magic  int32
	keyLen int32
	valLen int32
	Next   int32
}

type Cache struct {
	buf    []byte
	offset int32
}

func New(size int) *Cache {
	if size == 0 {
		size = tmpFileSize
	}

	buf := openTempFile(size)
	if buf == nil {
		return nil
	}

	c := &Cache{
		buf:    buf,
		offset: 0,
	}
	c.Clear()

	return c
}

func openTempFile(size int) []byte {
	path := os.Getenv("TOYSKKSERV_CACHE")
	if path == "" {
		path = filepath.Join(os.TempDir(), "toyskkserv.cache")
	}

	f, err := os.Create(path)
	if err != nil {
		slog.Error("Failed to open temporary file", "err", err)
		return nil
	}
	slog.Info("File cache", "path", path, "size", size)

	_, err = f.WriteAt([]byte{0}, int64(size-1))
	if err != nil {
		slog.Error("Failed to write temporary file", "err", err)
		f.Close()
		return nil
	}

	buf, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		slog.Error("Failed to map file to memory", "err", err)
		f.Close()
		return nil
	}

	return []byte(buf)
}

func (c *Cache) NewNode(key, val []byte) (node, int32) {
	if int(c.offset+nodeSize)+len(key)+len(val) >= len(c.buf) {
		return node{}, -1
	}

	addr := c.offset

	n := node{
		magic:  magicNum,
		keyLen: int32(len(key)),
		valLen: int32(len(val)),
		Next:   0,
	}
	c.offset += nodeSize

	copy(c.buf[c.offset:], key)
	c.offset += n.keyLen

	copy(c.buf[c.offset:], val)
	c.offset += n.valLen

	c.saveNode(n, addr)
	return n, addr
}

func (c *Cache) FindNode(addr int32) node {
	n := node{
		magic:  readInt32(c.buf[addr:]),
		keyLen: readInt32(c.buf[addr+4:]),
		valLen: readInt32(c.buf[addr+8:]),
		Next:   readInt32(c.buf[addr+12:]),
	}
	if n.magic != magicNum {
		panic(n.magic)
	}

	return n
}

func (c *Cache) Bytes(addr int32) ([]byte, []byte) {
	n := c.FindNode(addr)
	start := addr + nodeSize

	return c.buf[start : start+n.keyLen], c.buf[start+n.keyLen : start+n.keyLen+n.valLen]
}

func (c *Cache) Link(from int32, to int32) {
	n := c.FindNode(from)
	n.Next = to
	c.saveNode(n, from)
}

func (c *Cache) ReservedNode(key []byte) int32 {
	l := len(key)

	n := node{
		magic:  magicNum,
		keyLen: int32(l),
		valLen: int32(reservedSpace - l),
		Next:   0,
	}
	copy(c.buf[16:], key)
	c.saveNode(n, 0)

	return 0
}

func (c *Cache) Clear() {
	_ = c.ReservedNode([]byte{})
	c.offset = nodeSize + reservedSpace
}

func (c *Cache) saveNode(n node, addr int32) {
	writeInt32(c.buf[addr:], n.magic)
	writeInt32(c.buf[addr+4:], n.keyLen)
	writeInt32(c.buf[addr+8:], n.valLen)
	writeInt32(c.buf[addr+12:], n.Next)
}

func readInt32(buf []byte) int32 {
	return int32(buf[0]) + int32(buf[1])<<8 + int32(buf[2])<<16 + int32(buf[3])<<24
}

func writeInt32(buf []byte, n int32) {
	buf[0] = byte(n & int32(0xFF))
	buf[1] = byte(n >> 8)
	buf[2] = byte(n >> 16)
	buf[3] = byte(n >> 24)
}
