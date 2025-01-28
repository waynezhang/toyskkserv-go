package server

import (
	"bufio"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/waynezhang/eucjis2004decode/eucjis2004"
	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/defs"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
	"golang.org/x/text/transform"
)

type Server struct {
	listenAddr  string
	dictManager *dictionary.DictManager
	listener    net.Listener
}

func New(addr string, dm *dictionary.DictManager) *Server {
	return &Server{
		listenAddr:  addr,
		dictManager: dm,
	}
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

	switch req[0] {
	case defs.PROTOCOL_DISCONNECT:
		// CLIENT_END
		// Request to server: 0 + space + LF
		// Server terminates and disconnects after receiving the request
		slog.Info("Req type: disconnect")
		return false

	case defs.PROTOCOL_REQUEST:
		// CLIENT_REQUEST
		// Request to server: 1 + dictionary_key + space + LF
		// Answer if found: 1 + (/ + candidate) * (number of candidates) + / + LF
		// Answer if not found: 4 + dictionary_key + space + LF
		// The dictionary keys and candidates are all variable-length strings
		// The dictionary keys and candidates have the same character encoding
		// The primary encoding set of SKK is ASCII + euc-jp (note: UTF-8 can also be used in some implementations)
		slog.Info("Req type: request")
		s.dictManager.HandleRequest(req, w)
		w.Write([]byte{'\n'})
		return true

	case defs.PROTOCOL_VER:
		// CLIENT_VERSION
		// Request to server: 2 + space + LF
		// Answer: string including server version + space, e.g., dbskkd-cdb-2.00
		// Note: no known client parses this string
		// Implementation on dbskkd-cdb: returns the version string
		slog.Info("Req type: version")
		w.Write([]byte(defs.VersionString()))
		w.Write([]byte{' ', '\n'})
		return true

	case defs.PROTOCOL_HOST:
		// CLIENT_HOST
		// Request to server: 3 + space + LF
		// Answer: string including host information + space, e.g., localhost:127.0.0.1:
		// Note: no known client parses this string
		// Implementation on dbskkd-cdb: returns dummy string novalue:
		slog.Info("Req type: host")

		w.Write([]byte(s.listenAddr))
		w.Write([]byte{' ', '\n'})
		return true

	case defs.PROTOCOL_COMPLETION:
		// CLIENT_COMPLETION
		// Request to server: 4 + dictionary_key + space + LF
		// Same as CLIENT_REQUEST
		slog.Info("Req type: completion")
		s.dictManager.HandleCompletion(req, w)
		w.Write([]byte{'\n'})
		return true

	case 'c':
		// customized protocol
		slog.Info("Req type: customize command")
		return s.handleCustomizeCommand(req)

	default:
		slog.Error("Invalid request")
		return true
	}
}

func (s *Server) handleCustomizeCommand(req string) bool {
	key := strings.TrimSuffix(
		strings.TrimPrefix(req, string(defs.CUSTOMIZE_PROTOCOL)),
		" ",
	)
	switch key {
	case defs.CUSTOMIZE_PROTOCOL_RELOAD:
		urls := config.Shared().Dictionaries
		s.dictManager.DictionariesDidChange(urls)
		break
	}

	return true
}
