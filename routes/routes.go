package routes

import (
	"bytes"
	"fmt"
	"go-file-server/crud"
	"go-file-server/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func FileServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "POST":
		postHandler(w, r)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(*utils.Directory+r.URL.Path, "/")

	filename := r.URL.Query().Get("download")

	// Serve a file
	if filename != "" {
		buf := crud.DownloadDirectory(path, filename)

		w.Header().Set("Content-Disposition", "attachment; filename="+filename+".zip")
		w.Header().Set("Content-Type", "application/zip")

		http.ServeContent(w, r, filename+".zip", time.Now(), bytes.NewReader(buf.Bytes()))
	} else {

		// Ignore favicon.ico
		if !strings.Contains(path, "favicon.ico") {
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fileInfo, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			// If path is a directory, serve content as JSON
			if fileInfo.IsDir() {
				w.Header().Set("Content-Type", "application/json")
				files := crud.ReadDir(path)
				w.Write(files)
			} else {
				http.ServeFile(w, r, path)
			}
		}
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("File Upload Endpoint Hit")

	// 10 << 20 = 10 MiB
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("uploadF")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	// Write this byte array to our temporary file
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
