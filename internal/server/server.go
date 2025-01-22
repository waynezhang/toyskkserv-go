package server

import (
	"bufio"
	"log/slog"
	"net"
	"strings"

	"github.com/go-co-op/gocron/v2"
	"github.com/waynezhang/tskks/internal/config"
	"github.com/waynezhang/tskks/internal/dictionary"
	"golang.org/x/text/encoding/japanese"
)

type Server struct {
	listener    net.Listener
	dictManager *dictionary.DictManager
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start() {
	cfg := config.Shared()
	s.initializeDictionaries(cfg)
	s.startUpdateWatcher(cfg)

	addr, err := net.ResolveTCPAddr("tcp", cfg.ListenAddr)
	if err != nil {
		slog.Error("Failed to resolve addr", "addr", cfg.ListenAddr)
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

		go handleRequest(conn)
	}
}

func handleRequest(c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(c)
	decoder := japanese.EUCJP.NewDecoder()

	running := true
	for running {
		line, err := r.ReadString('\n')
		if err != nil {
			slog.Error("Failed to read from connection", "err", err)
			return
		}

		decoded, err := decoder.String(line)
		if err != nil {
			slog.Error("Failed to decode string", "req", line)
			c.Write([]byte("\n"))
			continue
		}

		req := strings.TrimSuffix(decoded, "\n")
		if len(req) == 0 {
			slog.Error("Empty reqeust")
			continue
		}
		slog.Info("Req received", "req", "["+req+"]", "cmd", req[0])

		switch req[0] {
		case '0':
			// CLIENT_END
			// Request to server: 0 + space + LF
			// Server terminates and disconnects after receiving the request
			slog.Info("Req type: disconnect")
			running = false
			break

		case '1':
			// CLIENT_REQUEST
			// Request to server: 1 + dictionary_key + space + LF
			// Answer if found: 1 + (/ + candidate) * (number of candidates) + / + LF
			// Answer if not found: 4 + dictionary_key + space + LF
			// The dictionary keys and candidates are all variable-length strings
			// The dictionary keys and candidates have the same character encoding
			// The primary encoding set of SKK is ASCII + euc-jp (note: UTF-8 can also be used in some implementations)
			slog.Info("Req type: request")
			res := dictionary.Shared().HandleRequest(req)
			slog.Info("Respnse", "res", "["+res+"]")
			c.Write([]byte(res + "\n"))
			break

		case '2':
			// CLIENT_VERSION
			// Request to server: 2 + space + LF
			// Answer: string including server version + space, e.g., dbskkd-cdb-2.00
			// Note: no known client parses this string
			// Implementation on dbskkd-cdb: returns the version string
			slog.Info("Req type: version")
			c.Write([]byte("tskks"))
			break

		case '3':
			// CLIENT_HOST
			// Request to server: 3 + space + LF
			// Answer: string including host information + space, e.g., localhost:127.0.0.1:
			// Note: no known client parses this string
			// Implementation on dbskkd-cdb: returns dummy string novalue:
			slog.Info("Req type: host")
			c.Write([]byte("\n"))
			break

		default:
			slog.Error("Invalid request")
			break
		}
	}
}

func (s *Server) initializeDictionaries(cfg *config.Config) {
	s.dictManager = dictionary.Shared()

	cfg.OnConfigChange(func() {
		s.dictManager.DictionariesDidChange()
	})
}

func (s *Server) startUpdateWatcher(cfg *config.Config) {
	cron := map[string]string{
		"daily":  "0 0 * * *",
		"weekly": "0 0 * * 1",
		"montly": "0 0 1 * *",
		"debug":  "* * * * *",
	}[cfg.UpdateSchedule]
	if cron == "" {
		slog.Info("Update schedule is disabled", "schedule", cfg.UpdateSchedule)
		return
	}
	slog.Info("Update schedule", "cron", cron)

	sh, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("Failed to start update watcher", "err", err)
		return
	}

	job := gocron.CronJob(cron, false)
	sh.NewJob(job, gocron.NewTask(func() {
		slog.Info("Checking updates")
		dictionary.UpdateDictionaries(cfg.Dictionaries, cfg.DictionaryDirectory, cfg.CacheDirectory)
	}))
	sh.Start()
}
