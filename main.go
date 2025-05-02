package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/oorrwullie/audioSlave/internal/app"
	"github.com/oorrwullie/audioSlave/internal/config"
)

func main() {
	configure := flag.Bool("config", false, "Reconfigure AudioSlave settings")
	flag.Parse()

	if *configure {
		if _, err := config.New(); err != nil {
			println("Configuration failed:", err.Error())
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx)
	if err != nil {
		os.Exit(1)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(application.Start)

	if err := eg.Wait(); err != nil && ctx.Err() == nil {
		application.Log.Error("AudioSlave exited with error", err)
		os.Exit(1)
	}

	application.Shutdown()
	application.Log.Info("AudioSlave shut down gracefully.")
}
