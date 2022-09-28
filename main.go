package main

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

func root(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(utils.RootPath+r.URL.Path, "/")

	if r.Method == http.MethodGet {

		filename := r.URL.Query().Get("download")
		upload := r.URL.Query().Get("upload")

		if filename != "" {
			buf := crud.DownloadDirectory(path, filename)

			w.Header().Set("Content-Disposition", "attachment; filename="+filename+".zip")
			w.Header().Set("Content-Type", "application/zip")

			http.ServeContent(w, r, filename+".zip", time.Now(), bytes.NewReader(buf.Bytes()))
		} else if upload != "" {
			fmt.Println("Upload endpoint!")
		} else {
			if !strings.Contains(path, "favicon.ico") {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}

				fileInfo, err := file.Stat()
				if err != nil {
					log.Fatal(err)
				}

				if fileInfo.IsDir() {
					w.Header().Set("Content-Type", "application/json")
					files := crud.ReadDir(path)
					w.Write(files)
				} else {
					http.ServeFile(w, r, path)
				}
			}
		}
	} else {
		log.Println("File Upload Endpoint Hit")

		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(32 << 20)
		// FormFile returns the first file for the given key `uploadFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, handler, err := r.FormFile("uploadF")
		if err != nil {
			log.Println("Error Retrieving the File")
			log.Println(err)
			return
		}
		defer file.Close()

		// m := r.MultipartForm
		// for _, v := range m.File {
		// 	fmt.Println("hey")
		// 	for _, f := range v {
		// 		file, err := f.Open()
		// 		if err != nil {
		// 			fmt.Println(err)
		// 			return
		// 		}
		// 		defer file.Close()
		// do something with the file data

		log.Printf("Uploaded File: %+v\n", handler.Filename)
		log.Printf("File Size: %+v\n", handler.Size)
		log.Printf("MIME Header: %+v\n", handler.Header)

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
		if err != nil {
			log.Println(err)
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		fmt.Fprintf(w, "Successfully Uploaded File\n")
		// 	}
		// }
	}

}

func main() {
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":2222", nil))
}
