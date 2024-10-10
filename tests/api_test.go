package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"orbcomm-ship-tracker/api"
	"orbcomm-ship-tracker/orbcomm"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func TestGetShipData(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a mock Orbcomm client
	mockOrbcommClient := &orbcomm.Client{
		BaseURL: "http://mock-orbcomm-api.com",
		APIKey:  "mock-api-key",
	}

	// Set up routes
	api.SetupRoutes(e, testDB, testS3Client, mockOrbcommClient, logger)

	// Create a request to pass to our handler
	req := httptest.NewRequest(http.MethodGet, "/ships/123456789", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/ships/:mmsi")
	c.SetParamNames("mmsi")
	c.SetParamValues("123456789")

	// Call the handler
	if err := api.GetShipData(c); err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	// Check the status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rec.Code)
	}

	// Check the response body
	var response []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Add more specific checks based on your expected response
	if len(response) == 0 {
		t.Errorf("Expected non-empty response; got empty")
	}
}

func TestGetAllShipsData(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a mock Orbcomm client
	mockOrbcommClient := &orbcomm.Client{
		BaseURL: "http://mock-orbcomm-api.com",
		APIKey:  "mock-api-key",
	}

	// Set up routes
	api.SetupRoutes(e, testDB, testS3Client, mockOrbcommClient, logger)

	// Create a request to pass to our handler
	req := httptest.NewRequest(http.MethodGet, "/ships", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	if err := api.GetAllShipsData(c); err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	// Check the status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rec.Code)
	}

	// Check the response body
	var response []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Add more specific checks based on your expected response
	if len(response) == 0 {
		t.Errorf("Expected non-empty response; got empty")
	}
}