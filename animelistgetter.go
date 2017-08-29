package main

import (
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

import "C"

var SeriesListRegexp = regexp.MustCompile(`<a href="/a/[a-zA-Z0-9-]+?" class="series-title" data-title="[^"]*?"(?: data-alt="[^"]*?")?>[^<]+`)
const SeriesListAppendix = `</a>`

//export ExtractSeriesList
func ExtractSeriesList(body string) ([]Series, error) {
	lines := strings.Split(body, "\n")
	series := make([]Series, 0)
	for _, linee := range lines {
		line := strings.TrimSpace(linee)
		if SeriesListRegexp.MatchString(line) {
			htmlLine := strings.NewReader(line + SeriesListAppendix)
			htmlElem, err := html.Parse(htmlLine)
			if err != nil {
				return nil, err
			}
			for (htmlElem.Type != html.ElementNode || htmlElem.Data != "a") && htmlElem.FirstChild != nil {
				htmlElem = htmlElem.FirstChild
				if htmlElem.Type == html.ElementNode && htmlElem.Data == "head" && htmlElem.NextSibling != nil {
					htmlElem = htmlElem.NextSibling
				}
			}
			if !(htmlElem.Type == html.ElementNode && htmlElem.Data == "a") {
				return nil, fmt.Errorf("parsing error: instead of <a> tag we found %#v <%v>", htmlElem.Type, htmlElem.Data)
			}
			var htmlHref string = ""
			var htmlAlt string = ""
			var htmlTitle string = ""
			for _, attr := range htmlElem.Attr {
				switch attr.Key {
				case "href":
					htmlHref = attr.Val
				case "data-alt":
					htmlAlt = attr.Val
				case "data-title":
					htmlTitle = attr.Val
				}
			}
			if htmlHref == "" || htmlTitle == "" {
				// silently ignore completely nonsense error
				continue
			}
			serie := Series{Title: htmlTitle, Slug: htmlHref[3:], Alt: htmlAlt}
			series = append(series, serie)
		}
	}
	return series, nil
}
