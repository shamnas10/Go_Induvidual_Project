package services

import (
	"database/sql"
	"datastream/logs"
)

func GetOutputFromClickHouse(query string, ClickhouseDb *sql.DB) (*sql.Rows, error) {
	rows, err := ClickhouseDb.Query(query)
	if err != nil {
		logs.Logger.Error("ERROR:", err)
		return nil, err
	}

	return rows, nil
}
