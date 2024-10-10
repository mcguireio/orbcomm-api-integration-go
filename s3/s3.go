package s3

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewS3Client creates a new S3 client
func NewS3Client() (*s3.S3, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

// ReadCSVFromS3 reads a CSV file from S3 and returns its contents
func ReadCSVFromS3(s3Client *s3.S3, bucket, key string) ([][]string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := s3Client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	reader := csv.NewReader(result.Body)
	return reader.ReadAll()
}

// WriteCSVToS3 writes CSV data to an S3 file
func WriteCSVToS3(s3Client *s3.S3, bucket, key string, data [][]string) error {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.WriteAll(data)
	if err != nil {
		return err
	}
	writer.Flush()

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	}

	_, err = s3Client.PutObject(input)
	return err
}

// GetShipListFromS3 retrieves the list of ships to track from an S3 CSV file
func GetShipListFromS3(s3Client *s3.S3, bucket, key string) ([]string, error) {
	csvData, err := ReadCSVFromS3(s3Client, bucket, key)
	if err != nil {
		return nil, err
	}

	var ships []string
	for i, row := range csvData {
		if i == 0 {
			continue // Skip header row
		}
		if len(row) > 0 {
			ships = append(ships, strings.TrimSpace(row[0]))
		}
	}
	return ships, nil
}

// ExportShipDataToS3 exports the latest ship data to an S3 CSV file
func ExportShipDataToS3(s3Client *s3.S3, bucket, key string, data []map[string]interface{}) error {
	csvData := [][]string{
		{"MMSI", "Name", "Latitude", "Longitude", "Speed", "Heading", "Timestamp"},
	}

	for _, ship := range data {
		row := []string{
			fmt.Sprintf("%v", ship["mmsi"]),
			fmt.Sprintf("%v", ship["name"]),
			fmt.Sprintf("%v", ship["latitude"]),
			fmt.Sprintf("%v", ship["longitude"]),
			fmt.Sprintf("%v", ship["speed"]),
			fmt.Sprintf("%v", ship["heading"]),
			fmt.Sprintf("%v", ship["timestamp"]),
		}
		csvData = append(csvData, row)
	}

	return WriteCSVToS3(s3Client, bucket, key, csvData)
}