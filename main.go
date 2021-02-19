package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const startURL = "https://en.wikipedia.org/wiki/Special:Random"
const targetURL = "https://en.wikipedia.org/wiki/Philosophy"
const wikipediaURL = "https://en.wikipedia.org"

var removeParentheses = regexp.MustCompile(`([a-zA-Z0-9]|\W)\(.*?.*?\)`)

func _findFirstLinkQuery(url string, history []string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	body := doc.Find("div#mw-content-text").Find("div.mw-parser-output")
	body.Find("table").Each(func(i int, t *goquery.Selection) {
		t.Remove()
	})

	var links []goquery.Selection
	body.Find("p").Each(func(i int, p *goquery.Selection) {
		paragraphHTML, err := p.Html()
		if err != nil {
			log.Fatal(err)
		}
		paragraphHTML = removeParentheses.ReplaceAllString(paragraphHTML, "")
		paragraphDoc, err := goquery.NewDocumentFromReader(strings.NewReader(paragraphHTML))
		if err != nil {
			log.Fatal(err)
		}
		paragraphDoc.Find("a").Each(func(i int, a *goquery.Selection) {
			links = append(links, *a)
		})
	})

	var firstLink string
	for _, link := range links {
		css, _ := link.Attr("class")
		href, _ := link.Attr("href")
		title, _ := link.Attr("title")
		if strings.Contains(css, "new") || strings.Contains(css, "external") || strings.HasPrefix(href, "#") || strings.HasPrefix(title, "Wikipedia:") || strings.HasPrefix(title, "Help:") || strings.HasPrefix(title, "File:") || strings.HasPrefix(href, "/wiki/File:") || strings.HasPrefix(href, "//upload.wikimedia.org/") || strings.HasPrefix(href, "https://en.wiktionary.org/") {
			continue
		}
		firstLink = wikipediaURL + href
		if _, visited := _find(history, firstLink); visited {
			continue
		}
		break
	}
	if len(firstLink) == 0 {
		return "", fmt.Errorf("Couldn't find first link in URL : %s", url)
	}

	return firstLink, nil
}

func main() {

	var links []string
	links = _findPhilosophy(startURL, links)

	fmt.Printf("Number of steps taken: %d\n", len(links))

}

func _findPhilosophy(url string, links []string) []string {

	link, err := _findFirstLinkQuery(url, links)
	fmt.Println(link)
	if err != nil {
		panic(err)
	}

	links = append(links, link)
	if link == targetURL {
		return links
	}

	return _findPhilosophy(link, links)
}

func _find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
