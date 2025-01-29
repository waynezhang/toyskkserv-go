package handler

import (
	"io"

	"github.com/waynezhang/toyskkserv/internal/dictionary"
)

type CompletionHandler struct {
	dm *dictionary.DictManager
}

func NewCompletionHandler(dm *dictionary.DictManager) *CompletionHandler {
	return &CompletionHandler{
		dm: dm,
	}
}

func (h CompletionHandler) Do(key string, w io.Writer) bool {
	// CLIENT_COMPLETION
	// Request to server: 4 + dictionary_key + space + LF
	// Same as CLIENT_REQUEST

	respWriter := newCandidateResponseWriter(w, key)
	h.dm.HandleCompletion(key, respWriter)
	respWriter.wrap()

	return true
}
