package orbcomm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents an Orbcomm API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// Ship represents the data structure for ship information
type Ship struct {
	MMSI      string    `json:"mmsi"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	Timestamp time.Time `json:"timestamp"`
}

// NewClient creates a new Orbcomm API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// GetShipData retrieves ship data for a given MMSI
func (c *Client) GetShipData(mmsi string) (*Ship, error) {
	url := fmt.Sprintf("%s/vessels/%s?api_key=%s", c.BaseURL, mmsi, c.APIKey)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var ship Ship
	err = json.NewDecoder(resp.Body).Decode(&ship)
	if err != nil {
		return nil, err
	}

	return &ship, nil
}

// ManageVesselList updates the list of vessels to track on the Orbcomm platform
func (c *Client) ManageVesselList(mmsiList []string) error {
	url := fmt.Sprintf("%s/vessel-list?api_key=%s", c.BaseURL, c.APIKey)
	
	payload := map[string]interface{}{
		"mmsi_list": mmsiList,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to manage vessel list. Status code: %d", resp.StatusCode)
	}

	return nil
}