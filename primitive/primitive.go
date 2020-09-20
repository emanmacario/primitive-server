package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode defines the shape used when transforming images
type Mode int

// Modes supported by the primitive package
const (
	ModeCombo          Mode = iota
	ModeTriangle       Mode = iota
	ModeRect           Mode = iota
	ModeEllipse        Mode = iota
	ModeCircle         Mode = iota
	ModeRotatedRect    Mode = iota
	ModeBeziers        Mode = iota
	ModeRotatedEllipse Mode = iota
	ModePolygon        Mode = iota
)

// WithMode is an option for the Transform function that will define the mode
// the user wants to use. By default, the ModeTriangle will be used
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

// Transform will take the provided image and apply a primitive transformation
// to it, then return a reader to the resulting image
func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	// Parse any additional (optional) primitive arguments
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}

	// Create temporary input file
	in, err := ioutil.TempFile("", fmt.Sprintf("in_*%s", ext))
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary input file")
	}
	defer os.Remove(in.Name())

	// Create temporary output file
	out, err := ioutil.TempFile("", fmt.Sprintf("out_*%s", ext))
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary output file")
	}
	defer os.Remove(out.Name())

	// Read image into in file
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("primitive: failed to copy image into temp input file")
	}

	// Run primitive w/ -i in.Name() -o out.Name()
	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, fmt.Errorf("primitive: failed to run primitive command. stdCombo=%s", stdCombo)
	}
	fmt.Println(stdCombo)

	// Read out into a reader, return reader, delete in and out
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("primitive: failed to copy output file into byte buffer")
	}
	return b, nil
}

// TODO: Refactor temp file creation into this function V
// func tempFile(prefix, ext string) (*os.File, error) {
// 	tmp, err := ioutil.TempFile("./images/", prefix)
// 	if err != nil {
// 		return nil, errors.New("primitive: failed to create temporary file")
// 	}
// 	defer os.remove(tmp.Name())
// 	return os.Create(fmt.Sprintf("%s%s", in.Name)), nil
// }

func primitive(inputFile, outputFile string, numShapes int, args ...string) (string, error) {
	argsStr := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes)
	args = append(strings.Fields(argsStr), args...)
	cmd := exec.Command("primitive", args...)
	stdoutStderr, err := cmd.CombinedOutput()
	return string(stdoutStderr), err
}
