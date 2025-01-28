package handler

import (
	"io"

	"github.com/waynezhang/toyskkserv/internal/defs"
)

type VersionHandler struct{}

func (VersionHandler) Do(req string, w io.Writer) bool {
	// CLIENT_VERSION
	// Request to server: 2 + space + LF
	// Answer: string including server version + space, e.g., dbskkd-cdb-2.00
	// Note: no known client parses this string
	// Implementation on dbskkd-cdb: returns the version string
	w.Write([]byte(defs.VersionString()))
	w.Write([]byte{' ', '\n'})
	return true
}
