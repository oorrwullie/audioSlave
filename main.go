package main

import (
	"context"
	"flag"
	"fmt"
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
		cfg, err := config.New()
		if err != nil {
			fmt.Printf("❌ Configuration failed: %v\n", err)
			os.Exit(1)
		}

		if cfg != nil {
			fmt.Println("✅ Configuration already exists.")
			fmt.Println("")
			fmt.Println("📌 To run the app automatically on startup, use:")
			fmt.Printf("   launchctl bootstrap gui/%d ~/Library/LaunchAgents/com.oorrwullie.audioSlave.plist\n", os.Getuid())
			fmt.Println("")
			fmt.Println("Or reboot your system — it will auto-start via launchd.")
		}

		os.Exit(0)
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
