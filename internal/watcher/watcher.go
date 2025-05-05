package watcher

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/oorrwullie/audioSlave/internal/config"
	"github.com/oorrwullie/audioSlave/internal/homebridge"
)

// Start launches the lockscreen-watcher binary and listens for LOCKED/UNLOCKED events.
func Start(ctx context.Context, cfg *config.Config) error {
	cmd := exec.CommandContext(ctx, "/usr/local/bin/lockscreen-watcher")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start lockscreen-watcher: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			return nil
		default:
			line := strings.TrimSpace(scanner.Text())
			switch strings.ToUpper(line) {
			case "LOCKED":
				log.Println("🔒 macOS Locked")
				log.Println("Turning OFF plug.")
				triggerPlug(cfg.GetClient(), cfg.PlugDevice, false)
			case "UNLOCKED":
				log.Println("🔓 macOS Unlocked")
				time.Sleep(2 * time.Second)
				if ok, _ := dacConnectedAndCorrectSampleRate(cfg); ok {
					log.Println("✅ DAC ready. Turning ON plug.")
					triggerPlug(cfg.GetClient(), cfg.PlugDevice, true)
				} else {
					log.Println("❌ DAC not ready. Not turning on plug.")
				}
			default:
				log.Printf("ℹ️ Unrecognized lockscreen-watcher output: %s", line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// dacConnectedAndCorrectSampleRate calls the dac-checker binary to validate DAC connection and sample rate.
func dacConnectedAndCorrectSampleRate(cfg *config.Config) (bool, error) {
	cmd := exec.Command("/usr/local/bin/dac-checker", cfg.DACName, cfg.SampleRate)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ dac-checker error: %v, output: %s", err, string(output))
		return false, nil
	}
	return strings.TrimSpace(string(output)) == "READY", nil
}

// triggerPlug toggles the smart plug via Homebridge.
func triggerPlug(hb *homebridge.Homebridge, device homebridge.Device, on bool) {
	if err := hb.TogglePlug(device, on); err != nil {
		log.Printf("❌ Failed to toggle plug: %v", err)
	} else {
		state := "off"
		if on {
			state = "on"
		}
		log.Printf("✅ Plug turned %s", state)
	}
}
