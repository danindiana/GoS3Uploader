package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/upload-multiple", uploadMultipleFiles)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	html := `<html>
<head>
  <title>File Upload</title>
</head>
<body>
  <form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="myFile" />
    <input type="submit" value="Upload Single File" />
  </form>
  <form enctype="multipart/form-data" action="/upload-multiple" method="post">
    <input type="file" name="myFiles" multiple />
    <input type="submit" value="Upload Multiple Files" />
  </form>
</body>
</html>`
	fmt.Fprint(w, html)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the form data
		r.ParseMultipartForm(10 << 20) // 10 MB max memory

		// Retrieve the file from form data
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		// Create a new file in the server's local storage
		dst, err := os.Create(handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the filesystem at the specified destination
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully Uploaded File\n")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func uploadMultipleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the form data
		r.ParseMultipartForm(10 << 20) // 10 MB max memory

		// Retrieve the file form data
		files := r.MultipartForm.File["myFiles"]

		for _, fileHeader := range files {
			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			// Create a new file in the server's local storage
			dst, err := os.Create(fileHeader.Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Copy the uploaded file to the filesystem at the specified destination
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		fmt.Fprintf(w, "Successfully Uploaded Multiple Files\n")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
