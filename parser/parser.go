package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Parser struct {
	StartUrl string
}

type ParseResult struct {
	Text     string    `json:"text"`
	Datetime time.Time `json:"datetime"`
}

func NewParser(startUrl string) *Parser {
	return &Parser{
		StartUrl: startUrl,
	}
}

func (p *Parser) Start() []ParseResult {
	var err error
	results := make([]ParseResult, 0)

	resultChan := make(chan string)

	reader, err := p.getUrlReader("https://eduface.ru/sites/list/region/2/type/1")
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
				countSites++
				go p.parseSite(url, resultChan)
			}
		})

		for countSites > 0 {
			value := <-resultChan
			//fmt.Printf("%d - %v\n", countSites, value)
			results = append(results, ParseResult{
				Text:     value,
				Datetime: time.Now(),
			})
			countSites--
		}
	}
	return results
}

func (p *Parser) getUrlReader(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (p *Parser) getUrl(url string) (string, error) {
	var err error
	var text []byte
	buf := make([]byte, 500)
	var countReaded int

	reader, err := p.getUrlReader(url)
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

func (p *Parser) parseSite(url string, resultChan chan string) {
	reader, err := p.getUrlReader(url)
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
			result = strings.ReplaceAll(strings.TrimSpace(h1Select.Text()), "\n", "")
		}
	}
}
