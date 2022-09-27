package main

import (
	"bytes"
	"go-file-server/crud"
	"go-file-server/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// var rootPath := os.Args[2]

func root(w http.ResponseWriter, r *http.Request) {
	// path := "." + strings.TrimSuffix(r.URL.Path, "/")
	path := strings.TrimSuffix(utils.RootPath+r.URL.Path, "/")

	keys, download := r.URL.Query()["download"]
	if download {
		filename := keys[0]
		// buf := crud.DownloadDirectory(r.URL.Path, filename)
		buf := crud.DownloadDirectory(path, filename)

		w.Header().Set("Content-Disposition", "attachment; filename="+filename+".zip")
		// w.Header().Set("Content-Disposition", "attachment; filename=prueba")
		w.Header().Set("Content-Type", "application/zip")

		http.ServeContent(w, r, filename+".zip", time.Now(), bytes.NewReader(buf.Bytes()))
		// http.ServeContent(w, r, "prueba", time.Now(), bytes.NewReader(buf.Bytes()))
	} else {

		// var path string
		// if r.URL.Path != "/" {
		// 	path = utils.RootPath + r.URL.Path
		// } else {
		// 	path = utils.RootPath
		// }
		// When opening a file in a new tab don't try to open favicon.ico
		// if path != "/favicon.ico" {
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

}

func main() {
	// fmt.Println(os.Args[2])
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":2222", nil))
}
