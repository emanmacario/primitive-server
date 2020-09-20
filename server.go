package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/emanmacario/primitive-server/primitive"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods(http.MethodGet)
	r.HandleFunc("/upload", handleUpload).Methods(http.MethodPost)

	prefix := "/images/"
	dir := http.Dir("./images/")
	fs := http.FileServer(dir)
	r.PathPrefix(prefix).Handler(http.StripPrefix(prefix, fs))

	log.Println("Server running on localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", r))

	inFile, err := os.Open("tmp/lenny.jpeg")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	out, err := primitive.Transform(inFile, "jpeg", 1)
	if err != nil {
		panic(err)
	}
	os.Remove("tmp/out.jpeg")
	outFile, err := os.Create("tmp/out.jpeg")
	if err != nil {
		panic(err)
	}
	io.Copy(outFile, out)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `<html><body>
		<form action="/upload" method="post" enctype="multipart/form-data">
			<input type="file" name="image">
			<button type="submit">Upload Image</button>
		</form>
		<body></html>`
	fmt.Fprint(w, html)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Retrieve file
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Set Content-Type header based on file extension
	ext := filepath.Ext(header.Filename)
	switch strings.ToLower(ext) {
	case ".jpg":
		fallthrough
	case ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		http.Error(w, "Invalid image type", http.StatusBadRequest)
		return
	}

	// Print image file metadata
	log.Printf("Uploaded file: %+v\n", header.Filename)
	log.Printf("File size: %+v\n", header.Size)
	log.Printf("MIME header: %+v\n", header.Header)
	log.Printf("Extension type: %s\n", ext)

	out, err := primitive.Transform(file, ext, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outFile, err := tempFile("", ext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	redirectURL := fmt.Sprintf("/%s", outFile.Name())
	http.Redirect(w, r, redirectURL, http.StatusFound)

}

func tempFile(prefix, ext string) (*os.File, error) {
	// Create temporary file for transformed image
	tmp, err := ioutil.TempFile("./images/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temporary file")
	}

	// Defer temporary file deletion
	defer func() {
		err := os.Remove(tmp.Name())
		if err != nil {
			log.Println(err.Error())
		}
	}()

	// Create a new file with temp file name and given extension
	return os.Create(fmt.Sprintf("%s%s", tmp.Name(), ext))
}
