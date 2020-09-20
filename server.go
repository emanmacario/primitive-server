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
	"text/template"

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
	// Retrieve image file from form
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse and validate file extension
	ext := filepath.Ext(header.Filename)
	extLower := strings.ToLower(ext)
	if extLower != ".jpg" && extLower != ".jpeg" && extLower != ".png" {
		http.Error(w, fmt.Sprintf("Unsupported filetype %s\n", extLower), http.StatusBadRequest)
		return
	}

	// Log uploaded image file metadata
	log.Printf("Uploaded file: %+v\n", header.Filename)
	log.Printf("File size: %+v\n", header.Size)
	log.Printf("MIME header: %+v\n", header.Header)
	log.Printf("Extension type: %s\n", ext)

	// Generate and return four different transformed images to the client
	a, err := genImage(file, ext, 33, primitive.ModeCombo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file.Seek(0, 0)
	b, err := genImage(file, ext, 33, primitive.ModeCircle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file.Seek(0, 0)
	c, err := genImage(file, ext, 33, primitive.ModeTriangle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file.Seek(0, 0)
	d, err := genImage(file, ext, 33, primitive.ModeRotatedRect)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
		{{range .}}.
			<img src="/{{.}}">
		{{end}}
		</body></html>`
	tmpl := template.Must(template.New("").Parse(html))
	log.Println([]string{a, b, d, c})
	err = tmpl.Execute(w, []string{a, b, c, d})
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
}

// Generates a primitive image and returns the image file name
func genImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(primitive.ModeCombo))
	if err != nil {
		return "", err
	}

	outFile, err := tempFile("", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)

	return outFile.Name(), nil
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
