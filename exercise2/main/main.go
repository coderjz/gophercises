package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	urlshort "github.com/coderjz/gophercises/exercise2"
)

func getFileContents(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s: %v", filepath, err)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v", filepath, err)
	}
	return b, nil
}

func main() {
	mux := defaultMux()

	var filePath string
	flag.StringVar(&filePath, "filepath", "urls.yaml", "location for the YAML file with the questions")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yaml, err := getFileContents(filePath)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	json := `{ "urls": [
		{ "path": "/gopher", "url": "www.gophercises.com" },
		{ "path": "/json-godoc", "url": "https://golang.org/pkg/encoding/json/" }
	] }`
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	err = urlshort.InitBoltDB()
	if err != nil {
		fmt.Println("Error reading from boltDB: ", err)
		fmt.Println("Starting the server on :8080 without bolt DB support")
		http.ListenAndServe(":8080", jsonHandler)
		return
	}

	boltDBHandler, err := urlshort.BoltDBHandler(jsonHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltDBHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
