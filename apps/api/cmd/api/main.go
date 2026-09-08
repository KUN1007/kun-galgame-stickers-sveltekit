package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kun-galgame-sticker-api/internal/app"
	"kun-galgame-sticker-api/pkg/config"
	"kun-galgame-sticker-api/pkg/logger"

	"github.com/joho/godotenv"
)

const shutdownTimeout = 15 * time.Second

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		port := os.Getenv("SERVER_PORT")
		if port == "" {
			port = "9421"
		}
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 2*time.Second)
		if err != nil {
			os.Exit(1)
		}
		_ = conn.Close()
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "error", err)
		os.Exit(1)
	}
	logger.Init(cfg.Server.Mode)

	application := app.New(cfg)
	defer application.Close()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdown
		slog.Info("shutting down")
		application.Close()
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := application.Fiber.ShutdownWithContext(ctx); err != nil {
			slog.Error("shutdown", "error", err)
		}
	}()

	addr := ":" + cfg.Server.Port
	slog.Info("listening", "addr", addr)
	if err := application.Fiber.Listen(addr); err != nil {
		slog.Error("listen", "error", err)
		os.Exit(1)
	}
}
