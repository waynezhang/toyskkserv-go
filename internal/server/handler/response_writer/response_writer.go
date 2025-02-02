package responsewriter

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

type Writer struct {
	io.Writer
	w                io.Writer
	hasResponseWrote bool
	key              string
	buf              *bytes.Buffer
	disposed         bool
}

var pool = sync.Pool{
	New: func() any {
		return &Writer{
			buf: bytes.NewBuffer(nil),
		}
	},
}

func New(w io.Writer, key string) *Writer {
	rw := pool.Get().(*Writer)

	rw.hasResponseWrote = false
	rw.key = key
	rw.w = w
	rw.buf.Reset()

	rw.disposed = false
	return rw
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if w.disposed {
		return 0, errors.New("Writer is not recycled properly")
	}
	n, err = 0, nil

	if !w.hasResponseWrote {
		w.hasResponseWrote = true
		_, err = w.buf.Write([]byte{'1'})
		if err != nil {
			return n, err
		}
	}

	n, err = w.buf.Write(p)
	return n + 1, err
}

func (w *Writer) Wrap() {
	if !w.hasResponseWrote {
		w.buf.Write([]byte{'4'})
		w.buf.Write([]byte(w.key))
		w.buf.Write([]byte{' ', '\n'})
	} else {
		w.buf.Write([]byte{'\n'})
	}
	w.w.Write(w.buf.Bytes())

	w.disposed = true
	pool.Put(w)
}
