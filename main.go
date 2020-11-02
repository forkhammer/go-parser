package main

import (
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var err error
	var text []byte
	buf := make([]byte, 500)
	var countReaded int

	response, err := http.Get("https://eduface.ru/sites/list/region/2/type/1")
	if err != nil {
		log.Fatalln(err)
	}

	for err == nil {
		countReaded, err = response.Body.Read(buf)
		text = append(text, buf[:countReaded]...)
	}
	if err != io.EOF {
		log.Fatalln(err)
	}
	doc, err := goquery.NewDocument(string(text))
}
