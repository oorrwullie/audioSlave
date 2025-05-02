package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/oorrwullie/audioSlave/internal/credentials"
	"github.com/oorrwullie/audioSlave/internal/homebridge"
)

type Config struct {
	DACName    string `json:"dac_name"`
	SampleRate string `json:"sample_rate"`
	BaseURL    string `json:"base_url"`
	PlugDevice string `json:"plug_device_id"`
	creds      *credentials.Credentials
	hbClient   *homebridge.Homebridge
}

const configPath = "/usr/local/etc/audioSlave/config.json"
const keychainService = "audioSlave"

func New() (*Config, error) {
	if _, err := os.Stat(configPath); err == nil {
		return loadFromDisk()
	}
	return createInteractive()
}

func loadFromDisk() (*Config, error) {
	c := &Config{}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}

	c.creds, err = credentials.New(keychainService)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	c.hbClient = homebridge.New(c.BaseURL, c.creds)
	return c, nil
}

func createInteractive() (*Config, error) {
	c := &Config{}
	reader := bufio.NewReader(os.Stdin)

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, err
	}

	fmt.Println("🔍 Scanning for audio devices...")
	out, err := exec.Command("system_profiler", "SPAudioDataType").Output()
	if err != nil {
		return nil, err
	}

	devices := parseAudioDevices(string(out))
	if len(devices) == 0 {
		return nil, fmt.Errorf("no audio devices found")
	}

	for i, device := range devices {
		fmt.Printf("[%d] %s\n", i, device)
	}
	defaultChoice := 1
	if defaultChoice >= len(devices) {
		defaultChoice = 0
	}
	fmt.Printf("Select your DAC device (default: %d): ", defaultChoice)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice := defaultChoice
	if input != "" {
		fmt.Sscanf(input, "%d", &choice)
	}
	if choice < 0 || choice >= len(devices) {
		return nil, fmt.Errorf("invalid choice")
	}
	c.DACName = devices[choice]

	fmt.Print("Enter desired sample rate [48000]: ")
	sampleRate, _ := reader.ReadString('\n')
	sampleRate = strings.TrimSpace(sampleRate)
	if sampleRate == "" {
		sampleRate = "48000"
	}
	c.SampleRate = sampleRate

	fmt.Print("Enter Homebridge base URL (e.g., http://homebridge.local:8581): ")
	baseURL, _ := reader.ReadString('\n')
	c.BaseURL = strings.TrimSpace(baseURL)

	if !strings.HasPrefix(c.BaseURL, "http://") && !strings.HasPrefix(c.BaseURL, "https://") {
		return nil, fmt.Errorf("invalid Homebridge base URL: must start with http:// or https://")
	}

	creds, err := credentials.New(keychainService)
	if err != nil {
		fmt.Print("Enter Homebridge UI username: ")
		u, _ := reader.ReadString('\n')
		u = strings.TrimSpace(u)

		fmt.Print("Enter Homebridge UI password (leave empty to auto-generate): ")
		p, _ := reader.ReadString('\n')
		p = strings.TrimSpace(p)

		if p == "" {
			p, err = credentials.GenerateRandomPassword(32)
			if err != nil {
				return nil, fmt.Errorf("failed to generate password: %w", err)
			}
			fmt.Printf("Generated secure random password: %s\n", p)
		}

		creds = &credentials.Credentials{ServiceName: keychainService}
		creds.SetUsername(u)
		creds.SetPassword(p)

		if err := creds.Save(); err != nil {
			return nil, fmt.Errorf("failed to save credentials: %w", err)
		}
	}
	c.creds = creds
	c.hbClient = homebridge.New(c.BaseURL, creds)

	// Discover devices
	devicesList, err := c.hbClient.ListDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to list Homebridge devices: %w", err)
	}

	if err := c.PromptPlugDeviceSelection(devicesList); err != nil {
		return nil, err
	}

	if err := c.save(); err != nil {
		return nil, err
	}

	fmt.Println("✅ Configuration saved!")
	return c, nil
}

func (c *Config) save() error {
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
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

// Getters
func (c *Config) GetBaseURL() string {
	return c.BaseURL
}

func (c *Config) GetUIUsername() string {
	return c.creds.GetUsername()
}

func (c *Config) GetUIPassword() string {
	return c.creds.GetPassword()
}

func (c *Config) GetPlugDeviceID() string {
	return c.PlugDevice
}

func (c *Config) GetClient() *homebridge.Homebridge {
	return c.hbClient
}

func (c *Config) GetCredentials() *credentials.Credentials {
	return c.creds
}

func (c *Config) PromptPlugDeviceSelection(devices []homebridge.Device) error {
	if len(devices) == 0 {
		return fmt.Errorf("no devices found to select from")
	}

	fmt.Println("Available plug devices:")
	for i, d := range devices {
		fmt.Printf("[%d] %s (%s)\n", i, d.Name, d.ID)
	}

	fmt.Print("Select device to control (enter number): ")
	var devIdx int
	fmt.Scanln(&devIdx)

	if devIdx < 0 || devIdx >= len(devices) {
		return fmt.Errorf("invalid device selection")
	}

	c.PlugDevice = devices[devIdx].ID
	return nil
}
