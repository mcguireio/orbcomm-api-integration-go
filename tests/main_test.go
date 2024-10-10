package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/lib/pq"
)

var (
	testDB       *sql.DB
	testS3Client *s3.S3
)

func TestMain(m *testing.M) {
	// Set up test environment
	var err error

	// Set up test database
	testDB, err = sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	// Set up test S3 client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Replace with your test region
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}
	testS3Client = s3.New(sess)

	// Run tests
	code := m.Run()

	// Tear down test environment
	if err := tearDown(); err != nil {
		log.Printf("Failed to tear down test environment: %v", err)
	}

	os.Exit(code)
}

// tearDown cleans up the test environment
func tearDown() error {
	// Clean up test database
	_, err := testDB.Exec("DROP TABLE IF EXISTS ship_data; DROP TABLE IF EXISTS ships;")
	if err != nil {
		return err
	}

	// Clean up test S3 bucket
	// Note: Be careful with this in a real test environment. You might want to use a dedicated test bucket.
	// _, err = testS3Client.DeleteBucket(&s3.DeleteBucketInput{
	// 	Bucket: aws.String(os.Getenv("TEST_S3_BUCKET")),
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}