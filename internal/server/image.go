package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// https://github.com/jayampathi-bac/ImageUploadApp/blob/main/main.go
// https://freshman.tech/file-upload-golang/

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) (string, error) {
	// Maximum upload of 20 MB files

	r.ParseMultipartForm(32 * MB)

	// r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	// if err := r.ParseMultipartForm(32 << 20); err != nil {
	// 	http.Error(w, "The uploaded file is too big. Please choose an file that's less than 20MB in size", http.StatusBadRequest)
	// 	return
	// }
	//take file from request
	file, handler, err := r.FormFile("image-post")
	if err != nil {
		if err.Error() == "http: no such file" {
			return "", nil
		}
		log.Println("Error Retrieving the File")
		log.Println(err)
		return "", err
	}
	log.Println("file read")

	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	if handler.Size > 20*MB {
		return "", fmt.Errorf("image is larger than 20MB")
	}

	// check file type
	// https://freshman.tech/file-upload-golang/
	buff := make([]byte, 512) //buff created because detect content type function works with max 512
	_, err = file.Read(buff)
	if err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" {
		return "", fmt.Errorf("The provided file format is not allowed")
	}

	_, err = file.Seek(0, io.SeekStart) // read buf changes the start point - we return it to the file's start
	if err != nil {
		return "", err
	}

	//create where to save file - we have chosen in the app (but can be saved in another server,git or special app, or directly as a bolb object in DB)
	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Println("Error Creating Folder for image")
		log.Println(err)
		return "", err
	}
	log.Println("file folder created")

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename)))
	if err != nil {
		log.Println("Error Creating the File for Image")
		return "", err
	}
	log.Println("file created")

	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("Error Copying Image")
		return "", err
	}

	log.Println("image upload successful")
	fullPath := dst.Name()
	Name := strings.TrimPrefix(fullPath, "./uploads/")
	return Name, nil
}
