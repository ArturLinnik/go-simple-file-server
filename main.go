package main

import (
	"bytes"
	"go-file-server/crud"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func root(w http.ResponseWriter, r *http.Request) {
	path := "." + strings.TrimSuffix(r.URL.Path, "/")

	directoryName, download := r.URL.Query()["download"]
	if download {
		buf := crud.DownloadDirectory(r.URL.Path, directoryName[0])

		w.Header().Set("Content-Disposition", "attachment; filename=test.zip")
		w.Header().Set("Content-Type", "application/zip")

		http.ServeContent(w, r, "test.zip", time.Now(), bytes.NewReader(buf.Bytes()))
	} else {
		// When opening a file in a new tab don't try to open favicon.ico
		if path != "./favicon.ico" {
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
				files := crud.ListDir(path)
				w.Write(files)
			} else {
				http.ServeFile(w, r, path)
			}

		}
	}

}

func main() {
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":2222", nil))
}
