package api

import (
	"datastream/config"
	"datastream/dataprocess"

	"datastream/logs"
	"errors"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UploadFileandInsertIntoKafka(w http.ResponseWriter, r *http.Request) {
	// Upload a csv file from an HTML form and get the file  and file header name.
	handler, file, err := UploadFile(w, r)
	if err != nil {
		return
	}
	MysqlDb, err := config.ConnectDatabase("mysql")
	if err != nil {
		logs.Logger.Error("Error connecting MySQL database: ", err)
	}

	go dataprocess.ConsumeContactActivitiesInKafka(MysqlDb)
	go dataprocess.ConsumeContactsInKafka(MysqlDb)

	// Load the query page templates.
	ExecuteQueryPageTemplates(w, r)

	// Copy the uploaded file to a directory and get the file path.
	filepath := CopyUploadedFile(handler, file)

	// Read and validate the  file, then send contact and contact activities to kafka
	dataprocess.ReadFile(filepath)

}

// execute homepage template
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/home/shamnas/Documents/test22/templates/HomePage.html")
	if err != nil {
		logs.Logger.Error("Error", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logs.Logger.Error("Error", err)

	}
}

var homePage = template.Must(template.ParseFiles("/home/shamnas/Documents/test22/templates/HomePage.html"))

// execute querypage template
func ExecuteQueryPageTemplates(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/home/shamnas/Documents/test22/templates/QueryPage.html")
	if err != nil {
		logs.Logger.Error("Error", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		logs.Logger.Error("Error", err)
	}
}

// Upload  file and check its type.
func UploadFile(w http.ResponseWriter, r *http.Request) (handler *multipart.FileHeader, file io.Reader, err error) {
	file, handler, err = r.FormFile("file")
	if err != nil {
		RenderUploadPageWithError(w, "Please choose a CSV file")
		return nil, nil, err
	}

	if !strings.HasSuffix(handler.Filename, ".csv") {
		RenderUploadPageWithError(w, "File is not a CSV")
		return nil, nil, errors.New("File is not a CSV")
	}

	return handler, file, nil
}

// Execute the homepage  and print an error message
func RenderUploadPageWithError(w http.ResponseWriter, errorMessage string) {
	data := struct {
		Error string
	}{Error: errorMessage}

	if err := homePage.Execute(w, data); err != nil {
		logs.Logger.Errorf(err.Error(), http.StatusInternalServerError)
	}
}

// Copy  uploaded file  and return the file path.
func CopyUploadedFile(handler *multipart.FileHeader, file io.Reader) string {
	// filename to store the uploaded file
	filename := filepath.Join("/home/shamnas/Documents/DatastreamFiles/", handler.Filename)
	out, err := os.Create(filename)
	if err != nil {
		logs.Logger.Error("Error", err)
	}
	defer out.Close()

	// Copy the file content to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		logs.Logger.Error("Error", err)
	}
	// Pass the path of the stored file to dataprocess.ReadFile
	return filename

}
