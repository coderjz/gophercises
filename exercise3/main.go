package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := 8080
	story, err := parseJSON()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	handler := NewHandler(story)
	fmt.Printf("Starting the server on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func parseJSON() (map[string]Arc, error) {
	var story map[string]Arc
	raw, err := ioutil.ReadFile("./gopher.json")
	if err != nil {
		return nil, fmt.Errorf("Error reading JSON file: %v", err)
	}

	err = json.Unmarshal(raw, &story)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON content: %v", err)
	}
	return story, nil
}

func NewHandler(story map[string]Arc) http.Handler {
	return handler{story}
}

type handler struct {
	story map[string]Arc
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}

	path = path[1:]

	arc, ok := h.story[path]
	if !ok {
		http.Error(w, "Arc not found", http.StatusNotFound)
		return
	}

	t := template.Must(template.ParseFiles("./tmpl/arc.html"))
	err := t.ExecuteTemplate(w, "arc.html", arc)
	if err != nil {
		log.Printf("Error executing html template : %v\n", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
	}
}

type Arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
