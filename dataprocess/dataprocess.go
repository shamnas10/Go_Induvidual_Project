package dataprocess

import (
	"database/sql"
	"datastream/config"
	"datastream/logs"
	"datastream/services"
	"datastream/types"
	"encoding/json"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

func ConsumeContactActivitiesInKafka(db *sql.DB) {
	err := godotenv.Load("/home/shamnas/Documents/test22/.env")
	if err != nil {
		logs.Logger.Error("Error loading .env file", err)
		return // Return to exit the function in case of an error
	}
	topic2 := os.Getenv("KAFKA_TOPIC2")
	table := os.Getenv("MYSQL_TABLE2")
	consumer, partitionConsumer, err := config.InitializeKafkaConsumer(topic2)
	if err != nil {
		logs.Logger.Error("Error initializing Kafka consumer", err)
		return // Handle the error gracefully by returning
	}
	defer consumer.Close()
	defer partitionConsumer.Close()

	for {
		message, err := services.ConsumeMessage(partitionConsumer)
		if err != nil {
			logs.Logger.Error("Error Consuming Messages From kafka", err)
		}
		if message != nil {
			var data []types.ContactActivity
			if err := json.Unmarshal([]byte(message), &data); err != nil {
				logs.Logger.Error("Error parsing Kafka message", err)
				continue
			}

			columns := []string{"ContactsID", "CampaignID", "ActivityType", "ActivityDate"}

			// Insert the data into the MySQL table
			count, err := services.InsertData(db, table, data, columns)
			if err != nil {
				logs.Logger.Error("Error inserting data into MySQL", err)
			} else {

				logs.Logger.Info("%d contact activity inserted", count)
			}

		}
	}
}
func ConsumeContactsInKafka(db *sql.DB) {
	err := godotenv.Load("/home/shamnas/Documents/test22/.env")
	if err != nil {
		logs.Logger.Error("Error loading .env file", err)
		return // Return to exit the function in case of an error
	}
	topic2 := os.Getenv("KAFKA_TOPIC1")
	table := os.Getenv("MYSQL_TABLE1")
	consumer, partitionConsumer, err := config.InitializeKafkaConsumer(topic2)
	if err != nil {
		logs.Logger.Error("Error initializing Kafka consumer", err)
		return // Handle the error gracefully by returning
	}
	defer consumer.Close()
	defer partitionConsumer.Close()
	var contactsbatch []types.Contacts

	for {

		message, err := services.ConsumeMessage(partitionConsumer)
		if err != nil {
			logs.Logger.Error("Error Consuming Messages From kafka", err)
		}
		if message != nil {
			var data types.Contacts
			if err := json.Unmarshal([]byte(message), &data); err != nil {
				logs.Logger.Error("Error parsing Kafka message", err)
				continue
			}
			contactsbatch = append(contactsbatch, data)

			columns := []string{"ID", "Name", "Email", "Details", "Status"}
			if len(contactsbatch) >= 100 {

				// Insert the data into the MySQL table
				count, err := services.InsertData(db, table, contactsbatch, columns)
				contactsbatch = nil
				if err != nil {
					logs.Logger.Error("Error inserting data into MySQL", err)
				} else {

					logs.Logger.Info("%d contacts inserted", count)
				}
			}

		}
	}
}

func GetResultFromClickHouse(w http.ResponseWriter, r *http.Request) {
	hidden := r.URL.Query().Get("hidden")
	// Parse and execute the "results.html" template with the result.
	tmpl, err := template.ParseFiles("/home/shamnas/Documents/test22/templates/ResultPage.html")
	if err != nil {
		logs.Logger.Error("ERROR:", err)
		return
	}
	var query string
	query = `SELECT
                    CampaignID,
                    opened
                    FROM contact_activity_last_three_month_Campaigns final
                    ORDER BY opened  LIMIT 5`
	if hidden == "Full Details" {
		query = `SELECT
                    CampaignID,
                    opened
                    FROM contact_activity_last_three_month_Campaigns final
                    ORDER BY opened  `
	}
	ClickhouseDb, err := config.ConnectDatabase("clickhouse") // Connect to Clickhouse
	rows, err := services.GetOutputFromClickHouse(query, ClickhouseDb)
	if err != nil {
		logs.Logger.Error("ERROR:", err)
		return
	}
	defer rows.Close()

	var results []types.ResultRow

	for rows.Next() {
		var row types.ResultRow
		err := rows.Scan(&row.CampaignId, &row.Count) // Adjust field names according to your data structure
		if err != nil {
			logs.Logger.Error("ERROR:", err)
			return
		}
		results = append(results, row)
	}

	err = tmpl.Execute(w, results)
	if err != nil {
		logs.Logger.Error("ERROR:", err)
		return
	}
}
