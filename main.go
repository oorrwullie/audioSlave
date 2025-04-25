package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	DACName          string `json:"dac_name"`
	SampleRate       string `json:"sample_rate"`
	HomebridgeOnURL  string `json:"homebridge_on_url"`
	HomebridgeOffURL string `json:"homebridge_off_url"`
}

var config Config

func main() {
	log.Println("Starting AudioSlave...")

	if len(os.Args) > 1 && os.Args[1] == "configure" {
		if err := configure(); err != nil {
			log.Fatalf("Configuration failed: %v", err)
		}
		return
	}

	if err := loadConfig("/usr/local/etc/audioSlave/config.json"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	go monitorSession()

	// Keep alive
	select {}
}

func loadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&config)
}

func monitorSession() {
	cmd := exec.Command("log", "stream", "--style", "syslog", "--predicate", "eventMessage CONTAINS \"Wake\" OR eventMessage CONTAINS \"Sleep\"")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start log stream: %v", err)
	}

	buf := make([]byte, 4096)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			log.Fatalf("Error reading log output: %v", err)
		}
		output := string(buf[:n])
		processLogOutput(output)
	}
}

func processLogOutput(output string) {
	output = strings.ToLower(output)
	if strings.Contains(output, "wake") {
		log.Println("Detected Unlock/Wake Event")
		time.Sleep(2 * time.Second) // Give system time to reconnect DAC
		if ok, _ := dacConnectedAndCorrectSampleRate(); ok {
			log.Println("DAC connected and correct sample rate. Turning ON plug.")
			triggerPlug(true)
		} else {
			log.Println("DAC not ready. Not turning on plug.")
		}
	} else if strings.Contains(output, "sleep") {
		log.Println("Detected Lock/Sleep Event")
		log.Println("Turning OFF plug.")
		triggerPlug(false)
	}
}

func dacConnectedAndCorrectSampleRate() (bool, error) {
	out, err := exec.Command("system_profiler", "SPAudioDataType").Output()
	if err != nil {
		return false, err
	}
	output := string(out)

	if !strings.Contains(output, config.DACName) || !strings.Contains(output, "Transport: USB") {
		return false, nil
	}

	re := regexp.MustCompile(`Current SampleRate:\s*(\\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return false, fmt.Errorf("could not find sample rate info")
	}

	if matches[1] == config.SampleRate {
		return true, nil
	}

	return false, nil
}

func triggerPlug(turnOn bool) {
	url := config.HomebridgeOffURL
	if turnOn {
		url = config.HomebridgeOnURL
	}

	reqBody := []byte(`{}`)
	resp, err := exec.Command("curl", "-X", "POST", url, "-H", "Content-Type: application/json", "-d", string(reqBody)).Output()
	if err != nil {
		log.Printf("Failed to toggle plug: %v", err)
		return
	}
	log.Printf("Plug toggle response: %s", string(resp))
}

func configure() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔍 Scanning for audio devices...")
	out, err := exec.Command("system_profiler", "SPAudioDataType").Output()
	if err != nil {
		return err
	}

	devices := parseAudioDevices(string(out))
	if len(devices) == 0 {
		return fmt.Errorf("no audio devices found")
	}

	for i, device := range devices {
		fmt.Printf("[%d] %s\n", i, device)
	}

	fmt.Print("Select your DAC device (enter number): ")
	var choice int
	fmt.Scanln(&choice)
	if choice < 0 || choice >= len(devices) {
		return fmt.Errorf("invalid choice")
	}
	config.DACName = devices[choice]

	fmt.Print("Enter desired sample rate [384000]: ")
	sampleRate, _ := reader.ReadString('\n')
	sampleRate = strings.TrimSpace(sampleRate)
	if sampleRate == "" {
		sampleRate = "384000"
	}
	config.SampleRate = sampleRate

	fmt.Print("Enter Homebridge base URL (e.g., http://homebridge.local:9000): ")
	hbURL, _ := reader.ReadString('\n')
	hbURL = strings.TrimSpace(hbURL)
	if hbURL == "" {
		return fmt.Errorf("homebridge URL is required")
	}
	config.HomebridgeOnURL = hbURL + "/plug/on"
	config.HomebridgeOffURL = hbURL + "/plug/off"

	if err := saveConfig("/usr/local/etc/audioSlave/config.json"); err != nil {
		return err
	}

	fmt.Println("✅ Configuration saved!")
	return nil
}

func parseAudioDevices(output string) []string {
	lines := strings.Split(output, "\n")
	var devices []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, ":") && !strings.HasPrefix(line, " ") {
			devices = append(devices, strings.TrimSuffix(line, ":"))
		}
	}
	return devices
}

func saveConfig(path string) error {
	os.MkdirAll("/usr/local/etc/audioSlave", 0755)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
