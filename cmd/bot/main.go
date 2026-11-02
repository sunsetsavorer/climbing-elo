package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sunsetsavorer/climbing-elo/internal/app"
	"github.com/sunsetsavorer/climbing-elo/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("bot stopped with error", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load(config.DefaultEnvFile)
	if err != nil {
		return err
	}

	bot, err := app.NewBot(cfg)
	if err != nil {
		return err
	}

	slog.Info("bot started")
	defer slog.Info("bot stopped")

	return bot.Run(ctx)
}
