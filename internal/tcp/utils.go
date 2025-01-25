package tcp

import (
	"log/slog"
	"net"

	"github.com/waynezhang/toyskkserv/internal/defs"
)

func SendReloadCommand(addr string) {
	sendTCPMessage(addr, string(defs.CUSTOMIZE_PROTOCOL)+defs.CUSTOMIZE_PROTOCOL_RELOAD)
}

func sendTCPMessage(addr string, message string) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		slog.Error("Failed to connect to server", "addr", addr, "err", err)
		return
	}
	defer c.Close()

	c.Write([]byte(message + "\n"))
}
