package services

import (
	"database/sql"
	"datastream/config"
	"datastream/logs"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetOutputFromOriginalClickHouse(t *testing.T) {
	ClickhouseDb, err := config.ConnectDatabase("clickhouse")
	if err != nil {
		logs.Logger.Error("clickhouse error", err)
	}
	// Test case 1: Valid Query
	t.Run("ValidQuery", func(t *testing.T) {
		query := `SELECT * FROM Contacts`
		rows, err := GetOutputFromClickHouse(query, ClickhouseDb)

		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
		if rows == nil {
			t.Error("Expected non-nil rows, but got nil rows.")
		}
	})

	// Test case 2: Invalid Query
	t.Run("InvalidQuery", func(t *testing.T) {
		query := "SELECT FROM Contacts" // Invalid query
		_, err := GetOutputFromClickHouse(query, ClickhouseDb)

		if err == nil {
			t.Error("Expected an error, but got no error.")
		}
	})

	// Test case 3: Connection Error
	t.Run("ConnectionError", func(t *testing.T) {
		query := `SELECT * FROM non_existent_table`
		_, err := GetOutputFromClickHouse(query, ClickhouseDb)

		if err == nil {
			t.Error("Expected an error due to connection failure or non-existent table, but got no error.")
		}
	})

	// Test case 4: Empty Query
	t.Run("EmptyQuery", func(t *testing.T) {
		query := ""
		_, err := GetOutputFromClickHouse(query, ClickhouseDb)

		if err == nil {
			t.Error("Expected an error for an empty query, but got no error.")
		}
	})

	// Test case 5: Nil Query
	t.Run("NilQuery", func(t *testing.T) {
		var query string // query is nil
		_, err := GetOutputFromClickHouse(query, ClickhouseDb)

		if err == nil {
			t.Error("Expected an error for a nil query, but got no error.")
		}
	})
}

type DatabaseConnector interface {
	ConnectDatabase(string) (*sql.DB, error)
}

type MockDatabaseConnector struct {
	mock.Mock
}

func (m *MockDatabaseConnector) ConnectDatabase(databaseName string) (*sql.DB, error) {
	args := m.Called(databaseName)
	return args.Get(0).(*sql.DB), args.Error(1)
}
func TestGetOutputFromClickHouse(t *testing.T) {
	// Create a real in-memory SQLite database for testing.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create a database connection: %v", err)
	}
	defer db.Close()

	// Create the required 'your_table' for the test in the SQLite database.
	_, err = db.Exec(`
        CREATE TABLE test_table (
            id INTEGER PRIMARY KEY,
            name TEXT
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create the table: %v", err)
	}
	// Insert sample data into the table.
	_, err = db.Exec(`INSERT INTO test_table (id,name) VALUES (1,"John"), (2,"Alice"), (3,"Bob")`)
	if err != nil {
		t.Fatalf("Failed to insert data into the table: %v", err)
	}
	// Define the expected query.
	expectedQuery := `SELECT * FROM test_table;`

	// Call the function with the query string and the database connection.
	rows, err := GetOutputFromClickHouse(expectedQuery, db)

	// Assert that the function returned the expected values.
	assert.NoError(t, err)
	assert.NotNil(t, rows)
}

func TestSendToKafkaModified(t *testing.T) {
	// Initialize the Kafka producer
	config.InitKafkaProducer()

	// Define a test topic for testing purposes
	testTopic := "test_topic"

	// Create a test message
	testMessage := []byte(`{"key": "value"}`)

	// Call the SendToKafka function
	SendToKafka(testTopic, testMessage, config.Producer)

	// Use a channel to wait for and validate the success message
	successChan := make(chan *sarama.ProducerMessage, 1)

	go func() {
		for success := range config.Producer.Successes() {
			successChan <- success
		}
	}()

	// Wait for the success message or a timeout
	select {
	case success := <-successChan:
		// Verify that the produced message is correct
		if success.Topic != testTopic {
			t.Errorf("Produced message topic does not match the expected topic.")
		}

	case <-time.After(5 * time.Second):
		t.Error("Timed out waiting for success message")
	}
}

func TestSendToKafkaWithInvalidTopic(t *testing.T) {
	config.InitKafkaProducer()

	// Define an invalid test topic for testing purposes
	invalidTestTopic := "invalid topic name"

	// Create a test JSON message
	testMessage := []byte(`{"key": "value"}`)

	// Call the SendToKafka function with the invalid test topic
	SendToKafka(invalidTestTopic, testMessage, config.Producer)

	// Use a channel to wait for and validate the success message
	successChan := make(chan *sarama.ProducerMessage, 1)
	go func() {
		for success := range config.Producer.Successes() {
			successChan <- success
		}
	}()

	// Wait for the success message or a timeout
	select {
	case success := <-successChan:
		// This block should not be reached when using an invalid topic.
		t.Errorf("Received a success message for an invalid topic: %s", success.Topic)
	case <-time.After(5 * time.Second):
		// This block should be reached due to the invalid topic.
		t.Logf("Test succeeded as expected: Timed out waiting for success message for an invalid topic")
	}
}
func BenchmarkSendToKafka(b *testing.B) {
	// Create a mock Kafka producer (you may need to set up a mock producer implementation)
	config.InitKafkaProducer()

	// Generate some sample JSON messages
	sampleMessages := []byte(`{"key1": "value1"}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Send messages to Kafka using the function and benchmark it
		SendToKafka("test-topic", sampleMessages, config.Producer)
	}
}
func TestConsumeMessage(t *testing.T) {
	config.InitKafkaProducer()
	defer config.CloseKafkaProducer()

	// Define your topic and message.
	topic := "test-topic"
	message := []byte("Hello Kafka!")

	// Send the message to Kafka using your SendToKafka function.
	SendToKafka(topic, message, config.Producer)
	consumer, partitionConsumer, err := config.InitializeKafkaConsumer(topic)
	if err != nil {
		logs.Logger.Error("Error initializing Kafka consumer", err)
		return // Handle the error gracefully by returning
	}
	defer consumer.Close()
	defer partitionConsumer.Close()
	// Wait for the message to be consumed or for a timeout.
	result, err := ConsumeMessage(partitionConsumer)
	if err != nil {
		t.Fatalf("Error consuming message: %v", err)
	}

	// Assert that the consumed message matches the sent message.
	assert.Equal(t, message, result)
}

func TestInsertData(t *testing.T) {
	// Open an SQLite in-memory database for testing.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table for inserting data.
	createTableQuery := "CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)"
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Define test data and columns.
	testData := []struct {
		ID   int
		Name string
	}{
		{1, "John"},
		{2, "Alice"},
	}

	columns := []string{"id", "name"}

	// Call the InsertData function to insert the test data.
	rowsAffected, err := InsertData(db, "test_table", testData, columns)
	if err != nil {
		t.Fatalf("InsertData failed: %v", err)
	}

	// Check that the expected number of rows were affected.
	expectedRowsAffected := int64(len(testData))
	if rowsAffected != expectedRowsAffected {
		t.Fatalf("Expected %d rows affected, but got %d", expectedRowsAffected, rowsAffected)
	}

	// Verify that the data was inserted correctly by querying the database.
	rows, err := db.Query("SELECT * FROM test_table")
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows and compare the data.
	var id int
	var name string
	count := 0
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			t.Fatalf("Error scanning row: %v", err)
		}
		expected := testData[count]
		if id != expected.ID || name != expected.Name {
			t.Fatalf("Row data does not match, expected: %v, got: (id=%d, name=%s)", expected, id, name)
		}
		count++
	}

	if count != len(testData) {
		t.Fatalf("Expected to retrieve %d rows, but got %d", len(testData), count)
	}
}

func TestInsertDataWithInvalidData(t *testing.T) {
	// Open an SQLite in-memory database for testing.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table for inserting data.
	createTableQuery := "CREATE TABLE test_table (ID INTEGER PRIMARY KEY, Name TEXT)"
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Attempt to insert invalid data (not a slice).
	invalidData := "not_a_slice"
	columns := []string{"ID", "Name"}

	_, err = InsertData(db, "test_table", invalidData, columns)
	if err == nil {
		t.Fatalf("Expected an error for invalid data, but got nil")
	}

	// Attempt to insert data with missing struct field.
	invalidData2 := []struct {
		ID   int
		Name string
	}{
		{1, "John"},
		{2, "Alice"},
	}

	invalidColumns := []string{"id", "age"} // "age" is not a field in the struct.

	_, err = InsertData(db, "test_table", invalidData2, invalidColumns)
	if err == nil {
		t.Fatalf("Expected an error for missing struct field, but got nil")
	}
}
