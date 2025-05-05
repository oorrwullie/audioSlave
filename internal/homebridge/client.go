package homebridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/oorrwullie/audioSlave/internal/credentials"
)

type Homebridge struct {
	baseUrl string
	creds   *credentials.Credentials
}

type Characteristic struct {
	AID         int      `json:"aid"`
	IID         int      `json:"iid"`
	UUID        string   `json:"uuid"`
	Type        string   `json:"type"`
	ServiceType string   `json:"serviceType"`
	ServiceName string   `json:"serviceName"`
	Description string   `json:"description"`
	Value       any      `json:"value"`
	Format      string   `json:"format"`
	Perms       []string `json:"perms"`
	CanRead     bool     `json:"canRead"`
	CanWrite    bool     `json:"canWrite"`
	EV          bool     `json:"ev"`
}

type AccessoryInformation struct {
	Manufacturer     string `json:"Manufacturer"`
	Model            string `json:"Model"`
	Name             string `json:"Name"`
	SerialNumber     string `json:"Serial Number"`
	FirmwareRevision string `json:"Firmware Revision"`
}

type Values struct {
	OutletInUse int `json:"OutletInUse"`
	On          int `json:"On"`
}

type Instance struct {
	Name                  string   `json:"name"`
	Username              string   `json:"username"`
	IPAddress             string   `json:"ipAddress"`
	Port                  int      `json:"port"`
	Services              []string `json:"services"`
	ConnectionFailedCount int      `json:"connectionFailedCount"`
	ConfigurationNumber   string   `json:"configurationNumber"`
}

type Device struct {
	AID                    int                  `json:"aid"`
	IID                    int                  `json:"iid"`
	UUID                   string               `json:"uuid"`
	Type                   string               `json:"type"`
	HumanType              string               `json:"humanType"`
	ServiceName            string               `json:"serviceName"`
	ServiceCharacteristics []Characteristic     `json:"serviceCharacteristics"`
	AccessoryInformation   AccessoryInformation `json:"accessoryInformation"`
	Values                 Values               `json:"values"`
	Instance               Instance             `json:"instance"`
	UniqueID               string               `json:"uniqueId"`
}

func New(baseUrl string, creds *credentials.Credentials) *Homebridge {
	return &Homebridge{
		baseUrl: baseUrl,
		creds:   creds,
	}
}

func (h *Homebridge) ListDevices() ([]Device, error) {
	token, err := h.getToken()
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	req, _ := http.NewRequest("GET", h.baseUrl+"/api/accessories", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)
		return nil, fmt.Errorf("failed to list devices (status %d): %s", resp.StatusCode, buf.String())
	}

	var devices []Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}

	return devices, nil
}

func (h *Homebridge) TogglePlug(device Device, on bool) error {
	token, err := h.getToken()
	if err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	url := fmt.Sprintf("%s/api/accessories/%s", h.baseUrl, device.UniqueID)

	payload := map[string]interface{}{
		"characteristicType": "On",
		"value":              on,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)
		return fmt.Errorf("homebridge API returned status %d: %s", resp.StatusCode, buf.String())
	}

	return nil
}

func (h *Homebridge) getToken() (string, error) {
	// Match curl payload exactly
	data := map[string]string{
		"username": h.creds.GetUsername(),
		"password": h.creds.GetPassword(),
	}
	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Match curl headers and method
	req, err := http.NewRequest("POST", h.baseUrl+"/api/auth/login", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send it
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// Optionally read response body for more detail
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)
		return "", fmt.Errorf("authentication failed (status %d): %s", resp.StatusCode, buf.String())
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func toInt(v interface{}) (int, bool) {
	f, ok := v.(float64)
	return int(f), ok
}
