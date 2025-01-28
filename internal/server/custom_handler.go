package server

import (
	"io"

	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/defs"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
)

type custom_handler struct {
	dm *dictionary.DictManager
}

func (h custom_handler) do(key string, w io.Writer) bool {
	switch key {
	case defs.CUSTOMIZE_PROTOCOL_RELOAD:
		urls := config.Shared().Dictionaries
		h.dm.DictionariesDidChange(urls)
		break
	}

	return true
}
