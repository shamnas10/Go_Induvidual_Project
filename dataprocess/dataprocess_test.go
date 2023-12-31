package dataprocess

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsValidJSON(t *testing.T) {
	t.Run("Valid JSON with Required Keys", func(t *testing.T) {
		validJSON := `{"country": "USA", "city": "New York", "dob": "1990-01-15"}`
		if !isValidJSON(validJSON) {
			t.Errorf("Expected valid JSON, but got invalid")
		}
	})

	t.Run("Valid JSON Missing One Required Key", func(t *testing.T) {
		invalidJSON := `{"country": "Canada", "dob": "1985-07-20"}`
		if isValidJSON(invalidJSON) {
			t.Errorf("Expected invalid JSON, but got valid")
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"name": "John", "age": 30}`
		if isValidJSON(invalidJSON) {
			t.Errorf("Expected invalid JSON, but got valid")
		}
	})

	t.Run("Empty JSON", func(t *testing.T) {
		emptyJSON := `{}`
		if isValidJSON(emptyJSON) {
			t.Errorf("Expected invalid JSON, but got valid")
		}
	})

	t.Run("using string", func(t *testing.T) {
		emptyJSON := `sdfgjzdf`
		if isValidJSON(emptyJSON) {
			t.Errorf("Expected invalid JSON, but got valid")
		}
	})

}

func TestIsValidString(t *testing.T) {
	t.Run("Valid String", func(t *testing.T) {
		validString := "hello, world"
		if !isValidString(validString) {
			t.Errorf("Expected valid string, but got invalid")
		}
	})

	t.Run("Invalid String", func(t *testing.T) {
		invalidString := "12345 hello"
		if isValidString(invalidString) {
			t.Errorf("Expected invalid string, but got valid")
		}
	})

	t.Run("Empty String", func(t *testing.T) {
		emptyString := " "
		if isValidString(emptyString) {
			t.Errorf("Expected invalid string, but got valid")
		}
	})

	t.Run("SpecialCharecters", func(t *testing.T) {
		SpecialCharecters := "@#"
		if isValidString(SpecialCharecters) {
			t.Errorf("Expected invalid string, but got valid")
		}
	})
}
func TestIsValidEmail(t *testing.T) {

	t.Run("Valid Email", func(t *testing.T) {
		validEmail := "johndoe@example.com"
		if !isValidEmail(validEmail) {
			t.Errorf("Expected valid email, but got invalid")
		}
	})

	t.Run("Invalid Email (Missing @)", func(t *testing.T) {
		invalidEmail := "johndoeexample.com"
		if isValidEmail(invalidEmail) {
			t.Errorf("Expected invalid email, but got valid")
		}
	})

	t.Run("Invalid Email (Invalid Characters)", func(t *testing.T) {
		invalidEmail := "john@doe@example.com"
		if isValidEmail(invalidEmail) {
			t.Errorf("Expected invalid email, but got valid")
		}
	})
	t.Run("Invalid Email (Invalid Characters)", func(t *testing.T) {
		invalidEmail := "john@doe@example.com"
		if isValidEmail(invalidEmail) {
			t.Errorf("Expected invalid email, but got valid")
		}
	})

	t.Run("string", func(t *testing.T) {
		string := "trgteas"
		if isValidEmail(string) {
			t.Errorf("Expected invalid email, but got valid")
		}
	})
}
func TestGenerateActivityType(t *testing.T) {

	// Run the function 1000 times to generate activity types.
	for i := 0; i < 1000; i++ {
		activityType := generateActivityType()

		// Check if the generated activity type is within the range of 1 to 7.
		if activityType < 1 || activityType > 7 {
			t.Errorf("Generated activity type out of range: %d", activityType)
		}
	}
}

func TestGetResultFromClickHouse(t *testing.T) {
	// Create a mock request for testing.
	requestBody := strings.NewReader("hidden=Full Details")
	req, err := http.NewRequest("POST", "/GetResult", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the HTTP response.
	rr := httptest.NewRecorder()

	// Call the function with the mock request and response recorder.
	GetResultFromClickHouse(rr, req)

	// Check the HTTP status code.
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Parse the response body (results) as a string.
	responseBody := rr.Body.String()

	// Define the expected data types.
	expectedCampaignIDType := "int"
	expectedOpenedType := "int"

	// Check if the response body contains the expected data types.
	if !strings.Contains(responseBody, expectedCampaignIDType) {
		t.Errorf("Expected response to contain data type: %s for CampaignID", expectedCampaignIDType)
	}

	if !strings.Contains(responseBody, expectedOpenedType) {
		t.Errorf("Expected response to contain data type: %s for opened", expectedOpenedType)
	}

}
