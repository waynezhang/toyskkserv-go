package server

import (
	"bufio"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/waynezhang/eucjis2004decode/eucjis2004"
	"github.com/waynezhang/toyskkserv/internal/defs"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
	"github.com/waynezhang/toyskkserv/internal/server/handler"
	"golang.org/x/text/transform"
)

type Server struct {
	listenAddr  string
	dictManager *dictionary.DictManager
	listener    net.Listener
	handlers    map[byte]requstHandler
}

func New(addr string, dm *dictionary.DictManager) *Server {
	s := &Server{
		listenAddr:  addr,
		dictManager: dm,
	}

	s.handlers = map[byte]requstHandler{}
	s.handlers[defs.PROTOCOL_DISCONNECT] = &handler.DisconnectHandler{}
	s.handlers[defs.PROTOCOL_REQUEST] = handler.NewCandidateHandler(dm)
	s.handlers[defs.PROTOCOL_VER] = &handler.VersionHandler{}
	s.handlers[defs.PROTOCOL_HOST] = handler.NewHostHandler(addr)
	s.handlers[defs.PROTOCOL_COMPLETION] = handler.NewCompletionHandler(dm)

	s.handlers[defs.CUSTOMIZE_PROTOCOL] = handler.NewCustomProtocolHandler(
		handler.DictManagerReload{Dm: dm},
	)

	return s
}

type requstHandler interface {
	Do(req string, w io.Writer) bool
}

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", s.listenAddr)
	if err != nil {
		slog.Error("Failed to resolve addr", "addr", s.listenAddr)
		panic(err)
	}

	slog.Info("Listen on", "addr", addr)
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		slog.Error("Failed to listen addr", "addr", addr)
		panic(err)
	}
	defer listener.Close()
	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept a connection", "err", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(transform.NewReader(c, eucjis2004.EUCJIS2004Decoder{}))

	running := true
	for running {
		line, err := r.ReadString('\n')
		if err != nil {
			slog.Info("Connect lost", "err", err)
			return
		}

		running = s.handleRequest(line, c)
	}
}

func (s *Server) handleRequest(req string, w io.Writer) bool {
	req = strings.TrimSuffix(req, "\n")
	if len(req) == 0 {
		slog.Error("Empty reqeust")
		return true
	}
	slog.Info("Req received", "req", "["+req+"]", "cmd", req[0])

	h := s.handlers[req[0]]
	if h == nil {
		slog.Error("Invalid request", "req", req[0])
		return true
	}

	return h.Do(strings.Trim(req[1:], " "), w)
}
