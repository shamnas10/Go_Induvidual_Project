package config

import (
	"database/sql"
	"datastream/logs"

	"errors"
	"fmt"
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDatabase(driverName string) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	switch driverName {
	case "mysql":
		db, err = loadMySQL()
		if err != nil {
			logs.Logger.Error("Error initializing MySQL database", err)
			return nil, err
		}
	case "clickhouse":
		db, err = loadClickhouse()
		if err != nil {
			logs.Logger.Error("Error initializing ClickHouse database", err)
			return nil, err
		}
	default:
		return nil, errors.New("unsupported database driver")
	}

	return db, nil
}

func loadMySQL() (*sql.DB, error) {
	if err := godotenv.Load("/home/shamnas/Documents/test24/.env"); err != nil {
		return nil, err
	}

	// Retrieve database connection details from environment variables
	dbUsername := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	// Construct the MySQL connection string
	connectionString := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

	// Initialize MySQL database connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}

// LoadClickhouse loads ClickHouse configuration from environment variables and returns a database connection.
func loadClickhouse() (*sql.DB, error) {
	if err := godotenv.Load("/home/shamnas/Documents/test24/.env"); err != nil {
		return nil, err
	}

	// Retrieve ClickHouse connection details from environment variables
	chHost := os.Getenv("CLICKHOUSE_HOST")
	chPort := os.Getenv("CLICKHOUSE_PORT")
	chUser := os.Getenv("CLICKHOUSE_USER")
	chPassword := os.Getenv("CLICKHOUSE_PASSWORD")
	chDatabase := os.Getenv("CLICKHOUSE_DATABASE")

	// Construct the ClickHouse connection string
	connectionString := fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s",
		chHost, chPort, chUser, chPassword, chDatabase)

	// Initialize ClickHouse database connection
	db, err := sql.Open("clickhouse", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
