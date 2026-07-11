package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/moges7624/gredis/internal/resp"
)

var CRLF = "\r\n"

type Config struct {
	Addr            string
	MaxConnections  int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		Addr:            ":6379",
		MaxConnections:  1000,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}
}

type Server struct {
	cfg    Config
	logger *slog.Logger

	sem chan struct{}
	wg  sync.WaitGroup
	nid atomic.Int64
	act atomic.Int64
}

func NewServer(cfg Config, logger *slog.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		sem:    make(chan struct{}, cfg.MaxConnections),
	}
}

func (s *Server) Run(ctx context.Context) error {
	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.cfg.Addr, err)
	}
	s.logger.Info("server listening", "addr", s.cfg.Addr, "max_connections", s.cfg.MaxConnections)

	go func() {
		<-ctx.Done()
		s.logger.Info("shutdown requested, closing listener")
		_ = ln.Close()
	}()

	s.acceptLoop(ctx, ln)

	return s.drain()
}

func (s *Server) acceptLoop(ctx context.Context, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				s.logger.Info("accept loop stopped")
				return
			}
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}
			s.logger.Error("accept failed", "error", err)
			continue
		}

		select {
		case s.sem <- struct{}{}:
		default:
			s.logger.Warn("connection limit reached, rejecting",
				"remote_addr", conn.RemoteAddr().String())
			_ = conn.Close()
			continue
		}

		id := s.nid.Add(1)
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer func() { <-s.sem }()
			s.handleConn(ctx, conn, id)
		}()
	}
}

func (s *Server) drain() error {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("graceful shutdown complete")
		return nil
	case <-time.After(s.cfg.ShutdownTimeout):
		s.logger.Warn("shutdown timeout exceeded, forcing exit",
			"active_connections", s.act.Load())
		return fmt.Errorf("shutdown timed out after %s with %d connections still active",
			s.cfg.ShutdownTimeout, s.act.Load())
	}
}

func (s *Server) handleConn(parent context.Context, conn net.Conn, id int64) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	defer conn.Close()

	s.act.Add(1)
	defer s.act.Add(-1)

	logger := s.logger.With(
		"conn_id", id,
		"remote_addr", conn.RemoteAddr().String(),
	)
	logger.Info("connection accepted")
	start := time.Now()
	defer func() {
		logger.Info("connection closed", "duration", time.Since(start).String())
	}()

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(s.cfg.IdleTimeout)); err != nil {
			logger.Error("set read deadline failed", "error", err)
			return
		}

		// TODO: reuse parser via a pool rather than create new one everytime
		parser := resp.NewParser(reader, logger)
		line, err := parser.Parse()
		if err != nil {
			switch {
			case errors.Is(err, io.EOF):
				logger.Info("client disconnected")
			case errors.Is(err, os.ErrDeadlineExceeded):
				logger.Info("connection idle timeout, closing")
			case parent.Err() != nil:
				logger.Info("connection closed for shutdown")
			default:
				logger.Warn("read error", "error", err)
			}
			return
		}

		if err := s.handleRequest(ctx, conn, logger, line); err != nil {
			logger.Warn("failed to handle line", "error", err)
			return
		}
	}
}

func (s *Server) handleRequest(
	ctx context.Context,
	conn net.Conn,
	_ *slog.Logger,
	val resp.Value,
) error {
	opCtx, cancel := context.WithTimeout(ctx, s.cfg.WriteTimeout)
	defer cancel()

	var resp string

	if err := conn.SetWriteDeadline(time.Now().Add(s.cfg.WriteTimeout)); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	cmd := strings.ToUpper(val.Array[0].Str)
	switch cmd {
	case "PING":
		resp = "+PONG\r\n"
	case "INFO":
		rep := "# Server\r\nredis_version:7.2.0\r\ntcp_port:6379\r\n"
		resp = ("$" + strconv.Itoa(len(rep)) + CRLF + string(rep) + CRLF)

	default:
		resp = "-ERR unknown command 'COMMAND'\r\n"
	}

	if _, err := conn.Write([]byte(resp)); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	select {
	case <-opCtx.Done():
		return opCtx.Err()
	default:
		return nil
	}
}
