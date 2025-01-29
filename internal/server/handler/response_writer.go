package handler

import (
	"bytes"
	"io"
)

type candidateResponseWriter struct {
	io.Writer
	w                io.Writer
	hasResponseWrote bool
	key              string
	buf              *bytes.Buffer
}

func newCandidateResponseWriter(w io.Writer, key string) *candidateResponseWriter {
	return &candidateResponseWriter{
		w:   w,
		key: key,
		buf: bytes.NewBuffer(nil),
	}
}

func (w *candidateResponseWriter) Write(p []byte) (n int, err error) {
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

func (w *candidateResponseWriter) wrap() {
	if !w.hasResponseWrote {
		w.buf.Write([]byte{'4'})
		w.buf.Write([]byte(w.key))
		w.buf.Write([]byte{' ', '\n'})
	} else {
		w.buf.Write([]byte{'\n'})
	}
	w.w.Write(w.buf.Bytes())
}
