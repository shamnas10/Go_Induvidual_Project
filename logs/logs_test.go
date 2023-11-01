package logs

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestSimpleLogger_Error(t *testing.T) {
	// Create a temporary log file for testing
	tempLogFile, err := ioutil.TempFile("", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempLogFile.Name())

	// Initialize the logger with the temporary log file
	logger, err := NewSimpleLogger(tempLogFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Call the Error method with a sample error message
	errorMessage := "Test error message"
	logger.Error(errorMessage, nil)

	// Close the logger to ensure the log message is flushed to the file
	_ = tempLogFile.Close()

	// Read the content of the temporary log file
	logContent, err := ioutil.ReadFile(tempLogFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the error message is present in the log content
	if !bytes.Contains(logContent, []byte(errorMessage)) {
		t.Errorf("Expected error message not found in log file")
	}
}

func TestSimpleLogger_Info(t *testing.T) {
	// Create a temporary log file for testing
	tempLogFile, err := ioutil.TempFile("", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempLogFile.Name())

	// Initialize the logger with the temporary log file
	logger, err := NewSimpleLogger(tempLogFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Call the Info method with a sample message
	infoMessage := "Test info message"
	logger.Info(infoMessage)

	// Close the logger to ensure the log message is flushed to the file
	_ = tempLogFile.Close()

	// Read the content of the temporary log file
	logContent, err := ioutil.ReadFile(tempLogFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the info message is present in the log content
	if !bytes.Contains(logContent, []byte(infoMessage)) {
		t.Errorf("Expected info message not found in log file")
	}
}

func TestSimpleLogger_Warning(t *testing.T) {
	// Create a temporary log file for testing
	tempLogFile, err := ioutil.TempFile("", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempLogFile.Name())

	// Initialize the logger with the temporary log file
	logger, err := NewSimpleLogger(tempLogFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Call the Warning method with a sample message
	warningMessage := "Test warning message"
	logger.Warning(warningMessage)

	// Close the logger to ensure the log message is flushed to the file
	_ = tempLogFile.Close()

	// Read the content of the temporary log file
	logContent, err := ioutil.ReadFile(tempLogFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the warning message is present in the log content
	if !bytes.Contains(logContent, []byte(warningMessage)) {
		t.Errorf("Expected warning message not found in log file")
	}
}
