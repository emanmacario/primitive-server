package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/emanmacario/primitive-server/primitive"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods(http.MethodGet)
	r.HandleFunc("/upload", handleUpload).Methods(http.MethodPost)

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
	switch ext {
	case ".jpg":
		fallthrough
	case ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		http.Error(w, "Invalid image type", http.StatusBadGateway)
		return
	}

	// Print image file metadata
	log.Printf("Uploaded file: %+v\n", header.Filename)
	log.Printf("File size: %+v\n", header.Size)
	log.Printf("MIME header: %+v\n", header.Header)
	log.Printf("Extension type: %s\n", ext)

	out, err := primitive.Transform(file, ext, 200)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.Copy(w, out)
}
