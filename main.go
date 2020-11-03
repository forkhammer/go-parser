package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

func getUrlReader(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func getUrl(url string) (string, error) {
	var err error
	var text []byte
	buf := make([]byte, 500)
	var countReaded int

	reader, err := getUrlReader(url)
	if reader != nil {
		defer reader.Close()
	}
	for err == nil {
		countReaded, err = reader.Read(buf)
		text = append(text, buf[:countReaded]...)
	}
	if err != io.EOF {
		return "", err
	}
	return string(text), nil
}

func parseSite(url string, resultChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	reader, err := getUrlReader(url)
	result := ""
	defer func() { resultChan <- result }()
	if err != nil {
		log.Println(err)
	} else {
		defer reader.Close()
		doc, _ := goquery.NewDocumentFromReader(reader)
		h1Select := doc.Find("h1")
		if h1Select.Length() > 0 {
			//fmt.Println(fmt.Sprintf("h1 = %v", h1Select.Text()))
			result = strings.TrimSpace(h1Select.Text())
		}
	}
}

func main() {
	var err error

	wg := sync.WaitGroup{}
	result := make(chan string)

	reader, err := getUrlReader("https://eduface.ru/sites/list/region/2/type/1")
	if err != nil {
		log.Println(err)
	} else {
		defer reader.Close()
		doc, _ := goquery.NewDocumentFromReader(reader)
		countSites := 0
		doc.Find("div.accordion-title .accordion-style4-wraplink a[href]").Each(func(index int, selection *goquery.Selection) {
			//fmt.Println(strings.TrimSpace(selection.Text()))
			url, exists := selection.Attr("href")
			if exists {
				//fmt.Println(url)
				wg.Add(1)
				countSites++
				go parseSite(url, result, &wg)
			}
		})

		for countSites > 0 {
			value := <-result
			fmt.Printf("%d - %v\n", countSites, value)
			countSites--
		}
	}
}
