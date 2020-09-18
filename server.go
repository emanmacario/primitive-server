package main

import (
	"io"
	"os"

	"github.com/emanmacario/primitive-server/primitive"
)

func main() {
	file, err := os.Open("tmp/lenny.jpeg")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	out, err := primitive.Transform(file, 5)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}
