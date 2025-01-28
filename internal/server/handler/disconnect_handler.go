package handler

import "io"

type DisconnectHandler struct{}

func (DisconnectHandler) Do(req string, w io.Writer) bool {
	// CLIENT_END
	// Request to server: 0 + space + LF
	// Server terminates and disconnects after receiving the request
	return false
}
