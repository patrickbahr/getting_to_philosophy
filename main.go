package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// const startURL = "https://en.wikipedia.org/wiki/Passenger_pigeon"

//edge case:
const startURL = "https://en.wikipedia.org/wiki/Geographic_coordinate_system"

//caso q deu pau
//https://en.wikipedia.org/wiki/List_of_Nobel_laureates_by_university_affiliation_II
const targetURL = "https://en.wikipedia.org/wiki/Philosophy"
const wikipediaURL = "https://en.wikipedia.org"

var removeParentheses = regexp.MustCompile(`([a-zA-Z0-9]|\W)\(.*?.*?\)`)

func findFirstLinkQuery(url string, history []string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	body := doc.Find("div#mw-content-text").Find("div.mw-parser-output")
	body.Find("table").Each(func(i int, t *goquery.Selection) {
		t.Remove()
	})

	var firstLink string
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
		paragraphDoc.Find("a").EachWithBreak(func(i int, link *goquery.Selection) bool {
			css, _ := link.Attr("class")
			href, _ := link.Attr("href")
			title, _ := link.Attr("title")
			//if !(strings.Contains(css, "new") || strings.Contains(css, "external") || strings.HasPrefix(href, "#") || strings.HasPrefix(title, "Wikipedia:") || strings.HasPrefix(title, "Help:") || strings.HasPrefix(title, "File:") || strings.HasPrefix(href, "/wiki/File:") || strings.HasPrefix(href, "//upload.wikimedia.org/") || strings.HasPrefix(href, "https://en.wiktionary.org/") || strings.HasPrefix(href, "https://es.wikipedia.org/wiki") || strings.HasPrefix(href, "https://en.wikisource.org/") || strings.HasPrefix(href, "https://commons.wikimedia.org") || strings.HasPrefix(href, "https://en.wikivoyage.org/wiki/") || strings.HasPrefix(href, "https://species.wikimedia.org/") || strings.HasPrefix(href, "https://ru.wikipedia.org/wiki/") || !(strings.HasPrefix(href, "/wiki/") || strings.HasPrefix(href, "https://en.wikipedia.org/wiki/"))) {
			validCSS := strings.Contains(css, "new") || strings.Contains(css, "external")
			validTitle := strings.HasPrefix(title, "Wikipedia:") || strings.HasPrefix(title, "Help:") || strings.HasPrefix(title, "File:")
			validHref := !((strings.HasPrefix(href, "/wiki/") && !strings.HasPrefix(href, "/wiki/File:") && !strings.HasPrefix(href, "/wiki/Template:") && !strings.HasPrefix(href, "/wiki/Module:")) || strings.HasPrefix(href, "https://en.wikipedia.org/wiki/"))
			if !(validCSS || validTitle || validHref) {
				f := wikipediaURL + href
				if _, visited := _find(history, f); !visited {
					firstLink = f
					return false
				}
			}
			return true
		})
	})

	if len(firstLink) == 0 {
		return "", fmt.Errorf("Couldn't find first link in URL : %s", url)
	}

	return firstLink, nil
}

func main() {

	totalSteps := 0
	count := 1
	for i := 0; i < count; i++ {
		var links []string
		links = findPhilosophy(startURL, links)
		steps := len(links)
		fmt.Printf("Number of steps taken: %d\n", steps)
		totalSteps += steps
	}

	fmt.Printf("Total number of steps taken: %d\n", totalSteps)
	fmt.Printf("Average number of steps taken: %d\n", (totalSteps / count))
	// link, _ := _findFirstLinkQuery(startUrl2, links)
	// fmt.Println(link)
}

const threshold = 1000

func findPhilosophy(url string, links []string) []string {

	link, err := findFirstLinkQuery(url, links)
	fmt.Println(link)
	if err != nil {
		panic(err)
	}

	if len(link) > threshold {
		panic("Threshold exceeded for this URL")
	}

	links = append(links, link)
	if link == targetURL {
		return links
	}

	return findPhilosophy(link, links)
}

func _find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
