package app

import (
	"context"

	"github.com/oorrwullie/audioSlave/internal/config"
	"github.com/oorrwullie/audioSlave/internal/logger"
	"github.com/oorrwullie/audioSlave/internal/watcher"
)

type App struct {
	Ctx context.Context
	Cfg *config.Config
	Log logger.Logger
}

func New(ctx context.Context) (*App, error) {
	log := logger.New()

	log.Info("🔧 Loading configuration...")
	cfg, err := config.New()
	if err != nil {
		log.Error("Failed to load config", err)
		return nil, err
	}

	return &App{
		Ctx: ctx,
		Cfg: cfg,
		Log: log,
	}, nil
}

func (a *App) Start() error {
	a.Log.Info("🚀 Starting AudioSlave...")
	return watcher.Start(a.Ctx, a.Cfg)
}

func (a *App) Shutdown() {
	a.Log.Info("🛑 Shutting down AudioSlave...")
}
