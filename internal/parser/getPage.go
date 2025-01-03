package parser

import (
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
)

func getPage(url, cookieValues string) (*html.Node, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", cookieValues)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func getPageFromFile(path string) (*html.Node, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file", err)
		return nil, err
	}
	defer f.Close()

	doc, err := html.Parse(f)

	if err != nil {
		return nil, err
	}

	return doc, nil
}
