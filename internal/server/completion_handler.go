package server

import (
	"io"

	"github.com/waynezhang/toyskkserv/internal/dictionary"
)

type completionHandler struct {
	dm *dictionary.DictManager
}

func (h completionHandler) do(key string, w io.Writer) bool {
	// CLIENT_COMPLETION
	// Request to server: 4 + dictionary_key + space + LF
	// Same as CLIENT_REQUEST
	respWriter := newCandidateResponseWriter(w, key)

	h.dm.HandleCompletion(key, respWriter)

	respWriter.close()
	return true
}
