package tasks

import (
	"database/sql"
	"time"

	"orbcomm-ship-tracker/database"
	"orbcomm-ship-tracker/orbcomm"
	"orbcomm-ship-tracker/s3"

	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

// StartOrbcommPoller starts the background task to poll the Orbcomm API
func StartOrbcommPoller(logger *zap.Logger, interval time.Duration, db *sql.DB, orbcommClient *orbcomm.Client, s3Client *s3.S3, s3Bucket string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		logger.Info("Polling Orbcomm API")

		// Get the list of ships to track from S3
		ships, err := s3.GetShipListFromS3(s3Client, s3Bucket, "ship_list.csv")
		if err != nil {
			logger.Error("Failed to get ship list from S3", zap.Error(err))
			continue
		}

		// Update the vessel list on Orbcomm platform
		err = orbcommClient.ManageVesselList(ships)
		if err != nil {
			logger.Error("Failed to manage vessel list", zap.Error(err))
			continue
		}

		// Fetch and store data for each ship
		var allShipData []map[string]interface{}
		for _, mmsi := range ships {
			shipData, err := orbcommClient.GetShipData(mmsi)
			if err != nil {
				logger.Error("Failed to get ship data", zap.String("mmsi", mmsi), zap.Error(err))
				continue
			}

			err = database.StoreShipData(db, shipData.MMSI, shipData.Name, shipData.Latitude, shipData.Longitude, shipData.Speed, shipData.Heading, shipData.Timestamp)
			if err != nil {
				logger.Error("Failed to store ship data", zap.String("mmsi", mmsi), zap.Error(err))
				continue
			}

			allShipData = append(allShipData, map[string]interface{}{
				"mmsi":      shipData.MMSI,
				"name":      shipData.Name,
				"latitude":  shipData.Latitude,
				"longitude": shipData.Longitude,
				"speed":     shipData.Speed,
				"heading":   shipData.Heading,
				"timestamp": shipData.Timestamp,
			})
		}

		// Export updated data to S3
		err = s3.ExportShipDataToS3(s3Client, s3Bucket, "latest_ship_data.csv", allShipData)
		if err != nil {
			logger.Error("Failed to export ship data to S3", zap.Error(err))
		}

		logger.Info("Finished polling Orbcomm API")
	}
}