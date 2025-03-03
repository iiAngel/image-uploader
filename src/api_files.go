package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type UploadFileResponse struct {
	Message          string `json:"message"`
	UploadedFileName string `json:"uploaded_file_name"`
}

var UploadedFiles []os.DirEntry

const charset = "abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"

func generateRandomString(length int) string {
	randomSeed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randomSeed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

func loadUploadedFiles() {
	var err error // uhmm lets ignore this
	UploadedFiles, err = os.ReadDir(LoadedConfig.FilesPath)
	if err != nil {
		log.Println(err)
	}
}

// NOTE: This takes a file name without an extension!
func searchFile(name string) string {
	if len(UploadedFiles) <= 0 {
		loadUploadedFiles()
	}

	for _, fileData := range UploadedFiles {
		extension := filepath.Ext(fileData.Name())
		filename, _ := strings.CutSuffix(fileData.Name(), extension)

		if filename == name {
			return filepath.Join(LoadedConfig.FilesPath, fileData.Name())
		}
	}

	return ""
}

func WriteUploadFileResponse(w http.ResponseWriter, response *UploadFileResponse, status int) {
	b, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func FilesUpload(w http.ResponseWriter, r *http.Request) {
	maxUploadSize := int64(LoadedConfig.MaxUploadSize) * 1024 * 1024

	err := r.ParseMultipartForm(0)
	if err != nil {
		log.Println(err)

		WriteUploadFileResponse(w, &UploadFileResponse{
			Message: "Unable to parse file",
		}, http.StatusBadRequest)
	}

	file, fileHeader, err := r.FormFile("uploaded_file")
	if err != nil {
		WriteUploadFileResponse(w, &UploadFileResponse{
			Message: "Unable to read file [Server]",
		}, http.StatusInternalServerError)
		return
	}

	defer file.Close()

	if fileHeader.Size > maxUploadSize {
		WriteUploadFileResponse(w, &UploadFileResponse{
			Message: "File is bigger than " + strconv.Itoa(LoadedConfig.MaxUploadSize) + "MBs",
		}, http.StatusBadRequest)
		return
	}

	randomName := generateRandomString(8)
	newFilePath := filepath.Join(LoadedConfig.FilesPath, randomName+filepath.Ext(fileHeader.Filename))

	createdFile, err := os.Create(newFilePath)
	if err != nil {
		log.Println(err)
		WriteUploadFileResponse(w, &UploadFileResponse{
			Message: "Unable to create file! [Server]",
		}, http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(createdFile, file); err != nil {
		log.Println(err)
		WriteUploadFileResponse(w, &UploadFileResponse{
			Message: "Unable to write to file! [Server]",
		}, http.StatusInternalServerError)
		return
	}

	// After all that stupid error handling, at this point this means
	// that the file has been saved!

	log.Println("Uploaded new file:", randomName)

	WriteUploadFileResponse(w, &UploadFileResponse{
		Message:          "File uploaded!",
		UploadedFileName: randomName,
	}, http.StatusOK)

	loadUploadedFiles()
}

func FilesShow(w http.ResponseWriter, r *http.Request) {
	requestQuery := r.URL.Query()
	fileName := requestQuery.Get("f")

	if fileName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filePath := searchFile(fileName)
	if filePath == "" {
		log.Println("Tried to get file", fileName, "but it wasn't found!")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filePath)))

	_, err = io.Copy(w, file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
