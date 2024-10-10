package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"orbcomm-ship-tracker/orbcomm"
)

func TestOrbcommClient(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"mmsi": "123456789",
			"name": "Test Ship",
			"latitude": 40.7128,
			"longitude": -74.0060,
			"speed": 10.5,
			"heading": 180,
			"timestamp": "2023-05-20T12:00:00Z"
		}`))
	}))
	defer server.Close()

	// Create an Orbcomm client with the mock server URL
	client := orbcomm.NewClient(server.URL, "test-api-key")

	// Test GetShipData
	ship, err := client.GetShipData("123456789")
	if err != nil {
		t.Fatalf("Failed to get ship data: %v", err)
	}

	// Check the returned data
	if ship.MMSI != "123456789" {
		t.Errorf("Expected MMSI 123456789, got %s", ship.MMSI)
	}
	if ship.Name != "Test Ship" {
		t.Errorf("Expected name 'Test Ship', got %s", ship.Name)
	}
	if ship.Latitude != 40.7128 {
		t.Errorf("Expected latitude 40.7128, got %f", ship.Latitude)
	}
	if ship.Longitude != -74.0060 {
		t.Errorf("Expected longitude -74.0060, got %f", ship.Longitude)
	}
	if ship.Speed != 10.5 {
		t.Errorf("Expected speed 10.5, got %f", ship.Speed)
	}
	if ship.Heading != 180 {
		t.Errorf("Expected heading 180, got %f", ship.Heading)
	}

	// Test ManageVesselList
	err = client.ManageVesselList([]string{"123456789", "987654321"})
	if err != nil {
		t.Fatalf("Failed to manage vessel list: %v", err)
	}
}