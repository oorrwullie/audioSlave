package watcher

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/oorrwullie/audioSlave/internal/config"
	"github.com/oorrwullie/audioSlave/internal/homebridge"
)

func Start(ctx context.Context, cfg *config.Config, hb *homebridge.Homebridge) error {
	cmd := exec.Command("log", "stream", "--style", "syslog", "--predicate", "eventMessage CONTAINS \"Wake\" OR eventMessage CONTAINS \"Sleep\"")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start log stream: %w", err)
	}

	buf := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			return nil
		default:
			n, err := stdout.Read(buf)
			if err != nil {
				return fmt.Errorf("error reading log output: %w", err)
			}
			output := string(buf[:n])
			processLogOutput(cfg, hb, output)
		}
	}
}

func processLogOutput(cfg *config.Config, hb *homebridge.Homebridge, output string) {
	output = strings.ToLower(output)
	if strings.Contains(output, "wake") {
		log.Println("Detected Unlock/Wake Event")
		time.Sleep(2 * time.Second)
		if ok, _ := dacConnectedAndCorrectSampleRate(cfg); ok {
			log.Println("DAC connected and correct sample rate. Turning ON plug.")
			triggerPlug(hb, cfg.GetPlugDeviceID(), true)
		} else {
			log.Println("DAC not ready. Not turning on plug.")
		}
	} else if strings.Contains(output, "sleep") {
		log.Println("Detected Lock/Sleep Event")
		log.Println("Turning OFF plug.")
		triggerPlug(hb, cfg.GetPlugDeviceID(), false)
	}
}

func dacConnectedAndCorrectSampleRate(cfg *config.Config) (bool, error) {
	out, err := exec.Command("system_profiler", "SPAudioDataType").Output()
	if err != nil {
		return false, err
	}
	output := string(out)

	if !strings.Contains(output, cfg.DACName) || !strings.Contains(output, "Transport: USB") {
		return false, nil
	}

	re := regexp.MustCompile(`Current SampleRate:\s*(\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return false, nil
	}

	return matches[1] == cfg.SampleRate, nil
}

func triggerPlug(hb *homebridge.Homebridge, deviceID string, on bool) {
	if err := hb.TogglePlug(deviceID, on); err != nil {
		log.Printf("Failed to toggle plug: %v", err)
	} else {
		state := "off"
		if on {
			state = "on"
		}
		log.Printf("Plug turned %s", state)
	}
}
