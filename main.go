package main

import (
	"encoding/json"
	"github.com/forkhammer/go-parser/parser"
	"log"
	"net/http"
)

func homeHandler(response http.ResponseWriter, request *http.Request) {
	newParser := parser.NewParser("https://eduface.ru/sites/list/region/2/type/1")
	data := newParser.Start()
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(data)
}

func main() {
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
