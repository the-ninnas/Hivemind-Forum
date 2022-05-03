package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// Checks if image filepath is allowed; returns "" if not allowed and returns the path if allowed.
func imageFileTypeValidator(filePath string) string {
	var count int
	allowedType := []string{".JPEG", ".JPG", ".PNG", ".GIF", ".jpeg", ".jpg", ".png", ".gif"}

	for _, v := range allowedType {
		if filePath == v {
			count++
		}
	}
	if count == 0 {
		return ""
	}
	return filePath
}

// Used in post.go
func uploadImageHandler(w http.ResponseWriter, r *http.Request) (string, string) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return "", ""
	}
	// Parse input
	err := r.ParseMultipartForm(20)
	if err != nil {
		fmt.Println("uploadimage parse error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", ""
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		// If no image, return empty strings
		if err.Error() == "http: no such file" {
			return "", ""
		}
		fmt.Println("uploadimage formfile error: ", err)
		tmpl.ExecuteTemplate(w, "400.html", nil)
		return "", ""
	}
	defer file.Close()

	filePath := filepath.Ext(handler.Filename)
	// Checks file type
	imageType := imageFileTypeValidator(filePath)
	if imageType == "" {
		msg := "Only JPEG, PNG and GIF images are allowed"
		fmt.Println(msg)
		return "", msg
	}

	// If file is more than 20mb, return error message
	if handler.Size > 20000000 {
		msg := "Image file size is over 20mb"
		fmt.Println(msg)
		return "", msg
	}

	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll("./images", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", ""
	}

	fileName := fmt.Sprintf("%d%v", time.Now().Unix(), filePath)

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./images/%s", fileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", ""
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", ""
	}

	fmt.Println("Image uploaded!")
	return fileName, ""
}
