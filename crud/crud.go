package crud

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"go-file-server/schemas"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ListDir(directory string) []byte {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	var response schemas.Response
	var filesArr []schemas.File
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		filesArr = append(filesArr, schemas.File{Name: fileInfo.Name(), IsDir: fileInfo.IsDir(), Size: fileInfo.Size(), ModTime: fileInfo.ModTime()})
	}

	response = schemas.Response{Files: filesArr, Path: directory}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}

func DownloadDirectory(URLPath string, directoryName string) *bytes.Buffer {
	URL := strings.TrimPrefix(URLPath, "/")
	directory := URL + "/" + directoryName

	fmt.Println(directory)

	file, err := os.Open(directory)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	if fileInfo.IsDir() {

		walker := func(path string, info os.FileInfo, err error) error {
			fmt.Printf("Crawling: %#v\n", path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Ensure that `path` is not absolute; it should not start with "/".
			// This snippet happens to work because I don't use
			// absolute paths, but ensure your real-world code
			// transforms path into a zip-root relative path.

			name := strings.TrimPrefix(path, URL)
			fmt.Println(name)
			f, err := writer.Create(name)
			if err != nil {
				return err
			}

			_, err = io.Copy(f, file)
			if err != nil {
				return err
			}

			return nil
		}
		err := filepath.Walk(directory, walker)
		if err != nil {
			panic(err)
		}

		writer.Close()

		return buf
	}

	return nil
}
