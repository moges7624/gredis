package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/moges7624/gredis/internal/command"
	"github.com/moges7624/gredis/internal/server"
	"github.com/moges7624/gredis/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := server.DefaultConfig()
	store := store.NewStore()
	dispatcher := command.NewDispatcher(store)

	server := server.NewServer(cfg, dispatcher.Handle, logger)

	if err := server.Run(ctx); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
