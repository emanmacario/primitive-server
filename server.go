package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/emanmacario/primitive-server/primitive"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods(http.MethodGet)
	r.HandleFunc("/upload", handleUpload).Methods(http.MethodPost)

	fmt.Println("Server running on localhost:5000")
	http.ListenAndServe(":5000", r)

	inFile, err := os.Open("tmp/lenny.jpeg")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	out, err := primitive.Transform(inFile, 1)
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
	// Parse input
	r.ParseMultipartForm(20)

	// Retrieve file
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error retrieving the file")
		log.Fatal(err)
		return
	}
	defer file.Close()

	// Print image file metadata
	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME header: %+v\n", handler.Header)

	// Write temporary file to server
	tempFile, err := ioutil.TempFile("images", "upload-*.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully uploaded image file\n")
}
