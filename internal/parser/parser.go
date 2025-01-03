package parser

import (
	"golang.org/x/net/html"
	"strings"
)

func findData(doc *html.Node, tag, substring string) []string {
	if doc != nil {
		var (
			crawler func(*html.Node)
			values  []string
		)

		crawler = func(node *html.Node) {
			if node.Type == html.ElementNode && node.Data == tag {
				if strings.Contains(node.FirstChild.Data, substring) {
					values = append(values, node.FirstChild.Data)
				}
			}
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				crawler(child)
			}
		}
		crawler(doc)

		return values
	}
	return nil
}
