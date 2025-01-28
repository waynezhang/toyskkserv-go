package handler

import "io"

type HostHandler struct {
	host string
}

func NewHostHandler(host string) *HostHandler {
	return &HostHandler{host: host}
}

func (h HostHandler) Do(req string, w io.Writer) bool {
	// CLIENT_HOST
	// Request to server: 3 + space + LF
	// Answer: string including host information + space, e.g., localhost:127.0.0.1:
	// Note: no known client parses this string
	// Implementation on dbskkd-cdb: returns dummy string novalue:
	w.Write([]byte(h.host))
	w.Write([]byte{' ', '\n'})
	return true
}
