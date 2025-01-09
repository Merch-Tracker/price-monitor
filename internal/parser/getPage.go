package parser

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

func getPage(url, cookieValues string) (*html.Node, error) {
	log.WithField("url length", len(url)).Debug("Fetching page")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithField("err", err).Error("Get Page | request")
		return nil, err
	}
	req.Header.Set("Cookie", cookieValues)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithField("err", err).Error("Get Page | response")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithField("code", resp.StatusCode).Error("Get Page | response status != 200")
		return nil, err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.WithField("err", err).Error("Get Page | html parse")
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
