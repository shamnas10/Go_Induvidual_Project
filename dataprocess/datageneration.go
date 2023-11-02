package dataprocess

import (
	cryptorand "crypto/rand"
	"regexp"
	"unicode"

	"datastream/config"
	"datastream/logs"
	"datastream/services"

	"datastream/types"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const (
	Active      = 1
	Bounce      = 2
	Unsubscribe = 5
	SpamAbuse   = 6
)

func generateRandomID() string {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		logs.Logger.Error("Error creating id", err)
	}

	randomBytes := make([]byte, 8)
	_, err = io.ReadFull(cryptorand.Reader, randomBytes) // Use rand.Reader from crypto/rand
	if err != nil {
		logs.Logger.Error("Error creating id", err)
	}

	randomString := fmt.Sprintf("%s-%x", uuidObj, randomBytes)
	randomString = strings.ReplaceAll(randomString, "-", "")

	return randomString
}

func generateActivityType() int {
	randomActivity := rand.Intn(1000) // Use rand.Intn from math/rand
	if randomActivity < 2 {
		randomActivity = 2
	} else if randomActivity < 3 {
		randomActivity = 6
	} else if randomActivity < 5 {
		randomActivity = 5
	} else if randomActivity <= 300 {
		randomActivity = 7
	} else if randomActivity <= 500 {
		randomActivity = 4
	} else if randomActivity <= 800 {
		randomActivity = 3
	} else {
		randomActivity = 1
	}
	return randomActivity
}
func generateRandomDate() time.Time {
	randomYear := 2023
	randomMonth := time.August
	randomDay := rand.Intn(28) + 1
	randomDate := time.Date(randomYear, randomMonth, randomDay, 0, 0, 0, 0, time.UTC)

	return randomDate
}
func isValidJSON(data string) bool {
	// Check if the length of the input string is greater than 3 characters
	if len(data) <= 3 {
		return false
	}
	var js map[string]interface{}
	if err := json.Unmarshal([]byte(data), &js); err != nil {
		return false
	}
	// Check if the JSON has the required keys and values
	requiredKeys := []string{"country", "city", "dob"}
	for _, key := range requiredKeys {
		if _, exists := js[key]; !exists {
			return false
		}
	}
	return true
}
func isValidString(s string) bool {
	firstChar := rune(s[0])
	return unicode.IsLetter(firstChar)
}

// isValidEmail checks if the provided email is in a valid format
func isValidEmail(email string) bool {
	// Regular expression pattern for a basic email format check
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(emailPattern, email)
	return match && err == nil
}
func ReadFile(filePath string) error {
	// Open the copied file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	count := 0
	lineNumber := 0

	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			logs.Logger.Info("File readed successfully")
			break
		}
		if err != nil {
			logs.Logger.Error("Error reading CSV", err)
			continue
		}
		if !isValidString(record[0]) {
			logs.Logger.ErrorinValidation("Invalid Type Name. It should be a string row:%d,value:%s", lineNumber, record[0])
			continue

		}
		// Check email format
		if !isValidEmail(record[1]) {
			logs.Logger.ErrorinValidation("Invalid Type Email Format . row:%d,value:%s", lineNumber, record[1])
			continue
		}
		// Check if the 3rd column is in JSON format
		if !isValidJSON(record[2]) {
			logs.Logger.ErrorinValidation("Invalid JSON format in the 3rd column. row:%d,value:%s", lineNumber, record[2])
			continue
		}

		// Check for whitespace
		if (record[0] == "") || (record[1] == "") || (record[2] == "") {
			logs.Logger.ErrorinValidation("FIle Conatin White Spaces.row:%d,name:%s", lineNumber, record[0])
			continue
		}

		id := generateRandomID()
		var contact types.Contacts
		contact.Name = record[0]
		contact.Email = record[1]
		contact.Details = record[2]
		contact.ID = id
		contact.Status = Active

		count++
		GenerateActivity(count, contact, id)
	}

	return nil
}
func GenerateActivity(count int, contact types.Contacts, id string) {
	err := godotenv.Load("/home/shamnas/Documents/test24/.env")
	if err != nil {
		logs.Logger.Error("Error loading .env file", err)
	}
	// Retrieve database connection details from environment variables
	topic1 := os.Getenv("KAFKA_TOPIC1")
	topic2 := os.Getenv("KAFKA_TOPIC2")
	var Activitybatch []types.ContactActivity

	for c := 1; c <= 100; c++ {

		actvt := generateActivityType()
		date := generateRandomDate()
		campaign := c
		dateStr := "2023-08-01"
		var temp int
		if actvt == 2 {
			temp = 2
			contact.Status = Bounce
		} else {
			temp = 1
		}

		for k := temp; k <= actvt; k++ {
			activityType := k
			var activityDate time.Time
			if k == 1 || k == 2 {

				activityDate, err = time.Parse("2006-01-02", dateStr)
				if err != nil {
					fmt.Println("Error parsing date:", err)
				}
			} else {
				activityDate = date
			}

			if k == 2 && actvt > 2 {
				continue
			} else if k == 5 && actvt == 7 {
				continue
			} else if k == 6 && actvt == 7 {
				continue
			}
			ca := types.ContactActivity{
				ContactsID:   id,
				CampaignID:   campaign,
				ActivityType: activityType,
				ActivityDate: activityDate,
			}

			if count%99 == 4 && k == 3 && c == 3 {
				Activitybatch = append(Activitybatch, ca)

			}
			if count%99 == 5 && k == 4 && c == 4 {
				Activitybatch = append(Activitybatch, ca)
			}

			Activitybatch = append(Activitybatch, ca)
			if k == 5 && actvt == 5 {
				contact.Status = Unsubscribe
			}
			if k == 6 && actvt == 6 {
				contact.Status = SpamAbuse
			}
			if contact.Status > 1 {
				break
			}

		}

		if contact.Status > 1 {
			break
		}

	}
	if len(Activitybatch) >= 1 {
		ActivityJSON, err := json.Marshal(Activitybatch)
		if err != nil {
			if err != nil {
				logs.Logger.Error("Error marshaling JSON", err)

			}
		}
		services.SendToKafka(topic2, ActivityJSON, config.Producer)
		Activitybatch = nil
	}

	DetailsJSON, err := json.Marshal(contact)
	if err != nil {
		if err != nil {
			logs.Logger.Error("Error marshaling JSON", err)
		}
	}
	services.SendToKafka(topic1, DetailsJSON, config.Producer)

}
