package crud

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"go-file-server/schemas"
	"go-file-server/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadDir(path string) []byte {
	fmt.Println(path)
	// Get files of a directory
	files, err := os.ReadDir(path)
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

	response = schemas.Response{Files: filesArr, Path: strings.TrimPrefix(path, utils.RootPath)}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}

func DownloadDirectory(path string, filename string) *bytes.Buffer {
	filePath := path + "/" + filename

	file, err := os.Open(filePath)
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

		walker := func(subPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(subPath)
			if err != nil {
				return err
			}
			defer file.Close()

			zipRoot := strings.TrimPrefix(subPath, path)
			f, err := writer.Create(zipRoot)
			if err != nil {
				return err
			}

			_, err = io.Copy(f, file)
			if err != nil {
				return err
			}

			return nil
		}
		err := filepath.Walk(filePath, walker)
		if err != nil {
			panic(err)
		}

		writer.Close()

		return buf
	}

	return nil
}
