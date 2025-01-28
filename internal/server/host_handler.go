package server

import "io"

type hostHandler struct {
	host string
}

func (h hostHandler) do(req string, w io.Writer) bool {
	// CLIENT_HOST
	// Request to server: 3 + space + LF
	// Answer: string including host information + space, e.g., localhost:127.0.0.1:
	// Note: no known client parses this string
	// Implementation on dbskkd-cdb: returns dummy string novalue:
	w.Write([]byte(h.host))
	w.Write([]byte{' ', '\n'})
	return true
}
