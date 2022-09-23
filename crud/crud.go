package crud

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"is_dir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type Response struct {
	Files []File `json:"files"`
	Path  string `json:"path"`
}

func ListDir(directory string) []byte {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	var filesArr []File
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		filesArr = append(filesArr, File{Name: fileInfo.Name(), IsDir: fileInfo.IsDir(), Size: fileInfo.Size(), ModTime: fileInfo.ModTime()})
	}

	response = Response{Files: filesArr, Path: directory}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}

func DownloadDirectory(URLPath string, directoryName string) *bytes.Buffer {
	// path := "." + strings.TrimSuffix(r.URL.Path, "/")
	// directory := strings.TrimPrefix(URLPath, "/")
	// fmt.Println(directory)
	// directory := "crud"
	URL := strings.TrimPrefix(URLPath, "/")
	directory := URL + "/" + directoryName

	// fmt.Println(URLPath)
	// fmt.Println(strings.TrimPrefix(URLPath, "/"))
	fmt.Println(directory)

	// file, err := os.Open(directory)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// fileInfo, err := file.Stat()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	// if fileInfo.IsDir() {

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

		// f, err := writer.Create(path)
		name := strings.TrimPrefix(path, URL)
		fmt.Println(name)
		// newName := directoryName + name
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
	// }

	// return nil
}
