package handler

import (
	"io"

	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/defs"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
)

type ReloadHandler interface {
	reload()
}

type CustomProtocolHandler struct {
	reloadHandler ReloadHandler
}

type DictManagerReload struct {
	Dm *dictionary.DictManager
}

func NewCustomProtocolHandler(reload ReloadHandler) *CustomProtocolHandler {
	return &CustomProtocolHandler{
		reloadHandler: reload,
	}
}

func (h CustomProtocolHandler) Do(key string, w io.Writer) bool {
	switch key {
	case defs.CUSTOMIZE_PROTOCOL_RELOAD:
		h.reloadHandler.reload()
		break
	}

	return true
}

func (h DictManagerReload) reload() {
	urls := config.Shared().Dictionaries
	h.Dm.DictionariesDidChange(urls)
}
