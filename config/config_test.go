package config

import (
	"os"
	"testing"
)

func TestConnectDatabase(t *testing.T) {
	t.Run("Test Connect to MySQL", func(t *testing.T) {
		db, err := ConnectDatabase("mysql")

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if db == nil {
			t.Error("Expected a non-nil database connection, but got nil")
		}
	})

	t.Run("Test Connect to ClickHouse", func(t *testing.T) {
		db, err := ConnectDatabase("clickhouse")

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if db == nil {
			t.Error("Expected a non-nil database connection, but got nil")
		}
	})

	t.Run("Test Unsupported Driver", func(t *testing.T) {
		db, err := ConnectDatabase("unsupported")

		if err == nil {
			t.Error("Expected an error, but got nil")
		}

		expectedErrorMsg := "unsupported database driver"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message: %s, but got: %s", expectedErrorMsg, err.Error())
		}

		if db != nil {
			t.Error("Expected a nil database connection, but got non-nil")
		}
	})
}

func TestInitKafkaProducer(t *testing.T) {
	// Set the KAFKA_BROKERS environment variable
	expectedBrokers := "localhost:9092"
	os.Setenv("KAFKA_BROKERS", expectedBrokers)
	defer os.Unsetenv("KAFKA_BROKERS")
	// Call the function you want to test
	InitKafkaProducer()
	if Producer == nil {
		t.Errorf("Producer should not be nil after initialization")
	}

	// Verify that the Kafka producer was initialized with the expected broker addresses
	actualBrokers := os.Getenv("KAFKA_BROKERS")
	if actualBrokers != expectedBrokers {
		t.Errorf("Broker address is incorrect. Got: %s, Expected: %s", actualBrokers, expectedBrokers)
	}
}
func TestCloseKafkaProducer(t *testing.T) {
	// Initialize the Kafka producer (assuming you have an InitKafkaProducer function)
	InitKafkaProducer()

	// Attempt to close the Kafka producer
	CloseKafkaProducer()

	// The function does not return an error, so no error handling is needed.
}
func TestInitializeKafkaConsumer(t *testing.T) {
	// Valid topic
	validTopic := "test_topic"
	consumer, partitionConsumer, err := InitializeKafkaConsumer(validTopic)
	if err != nil {
		t.Error("Expected no error  for a valid topic, but got error")
	}
	// Verify the returned values for a valid topic
	if consumer == nil {
		t.Error("Expected a non-nil Kafka consumer for a valid topic, but got nil")
	}

	if partitionConsumer == nil {
		t.Error("Expected a non-nil Kafka partition consumer for a valid topic, but got nil")
	}

	// Invalid topic
	invalidTopic := "invalid topic with spaces"
	invalidConsumer, invalidPartitionConsumer, err := InitializeKafkaConsumer(invalidTopic)
	if err == nil {
		t.Error("Expected  error  for a invalid topic, but got nill")
	}
	// Verify the returned values for an invalid topic
	if invalidConsumer != nil {
		t.Error("Expected a nil Kafka consumer for an invalid topic, but got a non-nil consumer")
	}

	if invalidPartitionConsumer != nil {
		t.Error("Expected a nil Kafka partition consumer for an invalid topic, but got a non-nil partition consumer")
	}
}
