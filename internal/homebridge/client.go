package homebridge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oorrwullie/audioSlave/internal/credentials"
)

type Homebridge struct {
	baseUrl string
	creds   *credentials.Credentials
}

func New(baseUrl string, creds *credentials.Credentials) *Homebridge {
	return &Homebridge{
		baseUrl: baseUrl,
		creds:   creds,
	}
}

func (h *Homebridge) TogglePlug(deviceID string, on bool) error {
	token, err := h.getToken()
	if err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	url := fmt.Sprintf("%s/api/accessories/%s", h.baseUrl, deviceID)
	payload := map[string]interface{}{
		"characteristicType": "On",
		"value":              on,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("homebridge API returned status %d", resp.StatusCode)
	}
	return nil
}

func (h *Homebridge) ListDevices() ([]Device, error) {
	token, err := h.getToken()
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	req, _ := http.NewRequest("GET", h.baseUrl+"/api/accessories", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accessories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("accessories request failed: %d", resp.StatusCode)
	}

	var accessories []accessory
	if err := json.NewDecoder(resp.Body).Decode(&accessories); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	var devices []Device
	for _, acc := range accessories {
		if acc.Category == 8 || acc.Category == 1 { // 8 = Outlet, 1 = Switch (some platforms mislabel)
			devices = append(devices, Device{
				ID:   acc.UUID,
				Name: acc.DisplayName,
			})
		}
	}

	return devices, nil
}

func (h *Homebridge) getToken() (string, error) {
	payload := map[string]string{
		"username": h.creds.GetUsername(),
		"password": h.creds.GetPassword(),
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(h.baseUrl+"/api/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("authentication failed")
	}

	var result authResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Token, nil
}

// Device represents a simplified accessory
type Device struct {
	ID   string
	Name string
}

// Homebridge returns these objects, simplified here
type accessory struct {
	UUID        string `json:"uuid"`
	DisplayName string `json:"displayName"`
	Category    int    `json:"category"`
}

type authResponse struct {
	Token string `json:"token"`
}
