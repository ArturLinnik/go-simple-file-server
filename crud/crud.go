package crud

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"go-file-server/schemas"
	"go-file-server/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadDir(path string) []byte {

	// Get files of a directory
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// Create an array with the files of the directory
	var response schemas.Response
	var filesArr []schemas.File
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		modTimeFormatted := fileInfo.ModTime().Format("02-01-2006 15:04:05")
		filesArr = append(filesArr, schemas.File{Name: fileInfo.Name(), IsDir: fileInfo.IsDir(), Size: fileInfo.Size(), ModTime: modTimeFormatted})
	}

	// Convert array to JSON and return it
	response = schemas.Response{Files: filesArr, Path: strings.TrimPrefix(path, *utils.Directory)}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}

func DownloadDirectory(path string, filename string) *bytes.Buffer {
	filePath := path + "/" + filename

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new buffer
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	if fileInfo.IsDir() {
		walker := func(subPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			// Open file
			file, err := os.Open(subPath)
			if err != nil {
				return err
			}
			defer file.Close()

			// Create zipped file
			zipRoot := strings.TrimPrefix(subPath, path)
			f, err := writer.Create(zipRoot)
			if err != nil {
				return err
			}

			// Copy file values into the zipped file
			_, err = io.Copy(f, file)
			if err != nil {
				return err
			}

			return nil
		}

		// Repeat this process recursively
		err := filepath.Walk(filePath, walker)
		if err != nil {
			panic(err)
		}

		writer.Close()
		return buf
	}

	return nil
}
