package routes

import (
	"bytes"
	"fmt"
	"go-file-server/crud"
	"go-file-server/utils"
	"io"
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
	// 10 << 20 = 10 MiB
	r.ParseMultipartForm(32 << 20)

	// Get file from form body
	file, handler, err := r.FormFile("upload-file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()

	// Get path from form body
	path := r.FormValue("path")

	// Create file
	dst, err := os.Create("./" + path + "/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
