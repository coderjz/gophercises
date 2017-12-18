package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coderjz/gophercises/exercise4/link"
)

func main() {
	filePaths := []string{"examples/ex1.html", "examples/ex2.html", "examples/ex3.html", "examples/ex4.html", "examples/ex5.html"}
	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		links, err := link.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n%+v\n\n", filePath, links)
	}
}
