package handler

import (
	"io"

	"github.com/waynezhang/toyskkserv/internal/dictionary"
	responsewriter "github.com/waynezhang/toyskkserv/internal/server/handler/response_writer"
)

type CandidateHandler struct {
	dm *dictionary.DictManager
}

func NewCandidateHandler(dm *dictionary.DictManager) *CandidateHandler {
	return &CandidateHandler{
		dm: dm,
	}
}

func (h CandidateHandler) Do(key string, w io.Writer) bool {
	// CLIENT_REQUEST
	// Request to server: 1 + dictionary_key + space + LF
	// Answer if found: 1 + (/ + candidate) * (number of candidates) + / + LF
	// Answer if not found: 4 + dictionary_key + space + LF
	// The dictionary keys and candidates are all variable-length strings
	// The dictionary keys and candidates have the same character encoding
	// The primary encoding set of SKK is ASCII + euc-jp (note: UTF-8 can also be used in some implementations)

	respWriter := responsewriter.New(w, key)
	h.dm.HandleRequest(key, respWriter)
	respWriter.Wrap()

	return true
}
