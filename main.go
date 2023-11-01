package main

import (
	"datastream/config"
	"datastream/routes"
	"net/http"
)

func main() {
	// Set up your routes
	routes.SetupRoutes()

	config.InitKafkaProducer()
	defer config.CloseKafkaProducer()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the HTTP server on port 8082
	http.ListenAndServe(":8081", nil)
}
