package services

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func InsertData(db *sql.DB, tableName string, data interface{}, columns []string) (int64, error) {
	// Check if data is a slice
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return 0, fmt.Errorf("Data must be a slice")
	}

	// Prepare placeholders for the columns
	numColumns := len(columns)
	placeholders := make([]string, numColumns)
	for i := range placeholders {
		placeholders[i] = "?"
	}

	// Construct the query with placeholders for columns
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var totalRowsAffected int64

	// Iterate over the data slice and insert rows
	for i := 0; i < val.Len(); i++ {
		// Extract values from the struct
		row := val.Index(i)
		values := make([]interface{}, numColumns)
		for j, column := range columns {
			field := row.FieldByName(column)
			if !field.IsValid() {
				return 0, fmt.Errorf("Field %s not found in the struct", column)
			}
			values[j] = field.Interface()
		}

		// Execute the query with data
		res, err := stmt.Exec(values...)
		if err != nil {
			return 0, err
		}

		// Get the number of rows affected
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}

		totalRowsAffected += rowsAffected
	}

	return totalRowsAffected, nil
}
