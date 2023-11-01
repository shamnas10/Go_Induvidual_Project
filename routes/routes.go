package routes

import (
	"datastream/api"
	"datastream/dataprocess"
	"net/http"
)

func SetupRoutes() {

	http.HandleFunc("/", api.HomePageHandler)
	http.HandleFunc("/UploadToKafka", api.UploadFileandInsertIntoKafka)
	http.HandleFunc("/GetResult", dataprocess.GetResultFromClickHouse)

}
