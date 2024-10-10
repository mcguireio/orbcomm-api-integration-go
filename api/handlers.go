package api

import (
	"database/sql"
	"net/http"
	"time"

	"orbcomm-ship-tracker/database"
	"orbcomm-ship-tracker/orbcomm"
	"orbcomm-ship-tracker/s3"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Handler holds the dependencies for API handlers
type Handler struct {
	DB            *sql.DB
	S3Client      *s3.S3
	OrbcommClient *orbcomm.Client
	Logger        *zap.Logger
}

// SetupRoutes sets up the API routes
func SetupRoutes(e *echo.Echo, db *sql.DB, s3Client *s3.S3, orbcommClient *orbcomm.Client, logger *zap.Logger) {
	h := &Handler{
		DB:            db,
		S3Client:      s3Client,
		OrbcommClient: orbcommClient,
		Logger:        logger,
	}

	e.GET("/ships/:mmsi", h.GetShipData)
	e.GET("/ships", h.GetAllShipsData)
}

// GetShipData handles the request to get data for a specific ship
func (h *Handler) GetShipData(c echo.Context) error {
	mmsi := c.Param("mmsi")
	data, err := database.GetShipData(h.DB, mmsi)
	if err != nil {
		h.Logger.Error("Failed to get ship data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get ship data"})
	}
	return c.JSON(http.StatusOK, data)
}

// GetAllShipsData handles the request to get data for all ships
func (h *Handler) GetAllShipsData(c echo.Context) error {
	// This is a simplified version. In a real-world scenario, you might want to implement pagination.
	rows, err := h.DB.Query(`
		SELECT DISTINCT ON (s.mmsi) s.mmsi, s.name, sd.latitude, sd.longitude, sd.speed, sd.heading, sd.timestamp
		FROM ships s
		JOIN ship_data sd ON s.id = sd.ship_id
		ORDER BY s.mmsi, sd.timestamp DESC
	`)
	if err != nil {
		h.Logger.Error("Failed to get all ships data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get all ships data"})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var mmsi, name string
		var lat, lon, speed, heading float64
		var timestamp time.Time
		err := rows.Scan(&mmsi, &name, &lat, &lon, &speed, &heading, &timestamp)
		if err != nil {
			h.Logger.Error("Failed to scan row", zap.Error(err))
			continue
		}
		results = append(results, map[string]interface{}{
			"mmsi":      mmsi,
			"name":      name,
			"latitude":  lat,
			"longitude": lon,
			"speed":     speed,
			"heading":   heading,
			"timestamp": timestamp,
		})
	}

	return c.JSON(http.StatusOK, results)
}