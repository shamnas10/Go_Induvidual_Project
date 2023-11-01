package api

import (
	"bytes"
	"io/ioutil"

	"io"
	"mime/multipart"
	"net/http"

	"net/http/httptest"
	"net/textproto"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestHomepageHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {

		t.Fatal(err)

	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(HomePageHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {

		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)

	}

}
func TestResultpageHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/GetResult", nil)

	if err != nil {

		t.Fatal(err)

	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(HomePageHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {

		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)

	}

}
func TestQuerytpageHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/UploadToKafka", nil)

	if err != nil {

		t.Fatal(err)

	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ExecuteQueryPageTemplates)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {

		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)

	}

}

func TestUploadFile(t *testing.T) {
	// Create a sample CSV file for testing
	csvContent := `shamnas	shamnas@gmail.com	{"country":"india","city":"kerala","dob":"29-09-2001"}
	// arun	arun@gmail.com	{"country":"usa","city":"washington","dob":"29-09-2002"}`

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	// Create a sample request with the CSV file
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Test case 1: Successful upload of a CSV file
	t.Run("SuccessfulUpload", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handler, file, err := UploadFile(recorder, req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if handler == nil {
			t.Errorf("Expected a file handler, got nil")
		}
		if file == nil {
			t.Errorf("Expected a file reader, got nil")
		}
		if !strings.HasSuffix(handler.Filename, ".csv") {
			t.Errorf("Expected file to have a .csv extension, got %s", handler.Filename)
		}
	})

	// Test case 2: Missing file in the request
	t.Run("MissingFile", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		req := httptest.NewRequest("POST", "/upload", body)          // Include an empty file
		req.Header.Set("Content-Type", writer.FormDataContentType()) // Set the correct content type
		recorder := httptest.NewRecorder()

		handler, file, err := UploadFile(recorder, req)

		if err == nil {
			t.Errorf("Expected error, but got no error: %v", err)
		}
		if handler != nil {
			t.Errorf("Expected a nil file handler, got %v", handler)
		}
		if file != nil {
			t.Errorf("Expected a nil file reader, got %v", file)
		}
	})

	// Test case 3: Uploaded file is not a CSV
	t.Run("NonCSVFile", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		// Create a sample request with a non-CSV file
		nonCSVContent := "This is not a CSV content"
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte(nonCSVContent))
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		_, _, err := UploadFile(recorder, req)

		if err == nil {
			t.Errorf("Expected an error, but got no error")
		}
		if !strings.Contains(err.Error(), "File is not a CSV") {
			t.Errorf("Expected error message to contain 'File is not a CSV', got %v", err)
		}
	})
}
func BenchmarkUploadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Create a test request with a sample file (similar to the test case)
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.csv")
		io.WriteString(part, "CSV File Content")
		writer.Close()

		req := httptest.NewRequest("POST", "http://example.com/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		_, _, _ = UploadFile(res, req)
	}
}

func TestCopyUploadedFile(t *testing.T) {
	// Create a temporary file to simulate an uploaded file
	tmpFile, err := os.CreateTemp("", "testfile.csv")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write some data to the temporary file
	data := []byte("Test CSV Data")
	_, err = tmpFile.Write(data)
	if err != nil {
		t.Fatalf("Failed to write data to the temporary file: %v", err)
	}

	// Create a multipart.FileHeader for the uploaded file
	fileHeader := &multipart.FileHeader{
		Filename: "testfile.csv",
		Header:   textproto.MIMEHeader{},
		Size:     int64(len(data)),
	}

	// Open the temporary file for reading
	tmpFile.Seek(0, 0)

	// Test the CopyUploadedFile function
	copiedFilePath := CopyUploadedFile(fileHeader, tmpFile)
	defer os.Remove(copiedFilePath) // Clean up copied file after the test

	// Check if the copied file exists
	if _, err := os.Stat(copiedFilePath); os.IsNotExist(err) {
		t.Errorf("The copied file does not exist at the expected path: %s", copiedFilePath)
	}

	// Open the copied file for reading
	copiedFile, err := os.Open(copiedFilePath)
	if err != nil {
		t.Fatalf("Failed to open the copied file: %v", err)
	}
	defer copiedFile.Close()

	// Read the content of the copied file
	copiedData, err := io.ReadAll(copiedFile)
	if err != nil {
		t.Fatalf("Failed to read the content of the copied file: %v", err)
	}

	// Compare the copied data with the original data
	if !bytes.Equal(copiedData, data) {
		t.Errorf("The copied file content does not match the original data")
	}
}
func BenchmarkCopyUploadedFile(b *testing.B) {
	// Create a temporary file to simulate the uploaded file
	tmpFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		b.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Generate some sample data
	sampleData := []byte("Sample file content")
	_, err = tmpFile.Write(sampleData)
	if err != nil {
		b.Fatalf("Error writing to temporary file: %v", err)
	}

	// Create a sample FileHeader
	handler := &multipart.FileHeader{
		Filename: "samplefile.txt",
	}

	// Open the temporary file for reading
	tmpFile.Seek(0, io.SeekStart)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Copy the file using the function and benchmark it
		_ = CopyUploadedFile(handler, tmpFile)

		// Optionally, you can check the result (file path) or perform additional benchmarks
	}
}
func TestRenderUploadPageWithError(t *testing.T) {
	// Helper function to perform the test.
	runTest := func(errorMessage string, expectError bool) {
		recorder := httptest.NewRecorder()
		RenderUploadPageWithError(recorder, errorMessage)

		// Check the response status code.
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
		}

		// Check if the response body contains or does not contain the error message as expected.
		if expectError {
			expectedErrorString := "Error: " + errorMessage
			if body := recorder.Body.String(); !strings.Contains(body, expectedErrorString) {
				t.Errorf("Expected response body to contain error message: %s", expectedErrorString)
			}
		} else {
			if body := recorder.Body.String(); strings.Contains(body, "Error:") {
				t.Errorf("Expected response body to not contain an error message, but it does.")
			}
		}
	}

	// Test with a non-empty error message.
	runTest("Test error message", true)

	// Test with an empty error message.
	runTest("", false)
}
