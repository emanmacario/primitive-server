package main

import (
	"io"
	"os"

	"github.com/emanmacario/primitive-server/primitive"
)

func main() {
	inFile, err := os.Open("tmp/lenny.jpeg")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	out, err := primitive.Transform(inFile, 100)
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
