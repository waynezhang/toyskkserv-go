package server

import "io"

type candidateResponseWriter struct {
	io.Writer
	w                io.Writer
	hasResponseWrote bool
	key              string
}

func newCandidateResponseWriter(w io.Writer, key string) *candidateResponseWriter {
	return &candidateResponseWriter{w: w, key: key}
}

func (w *candidateResponseWriter) Write(p []byte) (n int, err error) {
	n, err = 0, nil

	if !w.hasResponseWrote {
		w.hasResponseWrote = true
		n, err = w.Write([]byte{'1'})
		if err != nil {
			return n, err
		}
	}

	wrote, err := w.w.Write(p)
	return n + wrote, err
}

func (w *candidateResponseWriter) close() {
	if !w.hasResponseWrote {
		w.w.Write([]byte{'4'})
		w.w.Write([]byte(w.key))
		w.w.Write([]byte{' ', '\n'})
	} else {
		w.w.Write([]byte{'\n'})
	}
}
