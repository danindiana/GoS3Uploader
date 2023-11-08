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
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	dst, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	counter := &WriteCounter{}
	if _, err = io.Copy(dst, io.TeeReader(file, counter)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func uploadMultipleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// The request's ParseMultipartForm method must be called before accessing file data
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // Adjust the size as needed
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	files := r.MultipartForm.File["myFiles"]
	for i, fileHeader := range files {
		fmt.Printf("Uploading file %d/%d: %s\n", i+1, len(files), fileHeader.Filename)

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		dst, err := os.Create(fileHeader.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		counter := &WriteCounter{}
		if _, err = io.Copy(dst, io.TeeReader(file, counter)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Finished uploading file %d/%d: %s\n", i+1, len(files), fileHeader.Filename)
	}

	fmt.Fprintf(w, "Successfully Uploaded Multiple Files\n")
}


type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", bytesToHuman(wc.Total))
}

func bytesToHuman(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
