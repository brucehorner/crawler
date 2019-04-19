package main

import (
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
)

func extractHyperlinks(thisURL string, tokenizer *html.Tokenizer, domainsticky bool) []string {
	var URLs []string
	if tokenizer == nil || len(thisURL) == 0 {
		return URLs
	}
L:
	for {
		tokenType := tokenizer.Next()
		switch {
		case tokenType == html.ErrorToken:
			break L
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			// check for robots
			if token.Data == "meta" {
				name := token.Attr[0]
				if name.Key == "name" && name.Val == "robots" {
					content := token.Attr[0]
					if content.Key == "content" {
						if strings.Contains(content.Val, "nofollow") || strings.Contains(content.Val, "noindex") {
							log.Println("robots ignore directive:", content.Val, "for", thisURL)
							return URLs
						}
					}
				}
			}

			// skip every tag except anchor tags
			if token.Data == "a" {
				for _, a := range token.Attr {
					if a.Key == "href" {
						childURL := strings.Trim(a.Val, " \t\n")
						if len(childURL) == 0 {
							continue
						}

						uThis, _ := url.Parse(thisURL)
						uChild, e := url.Parse(childURL)
						if e != nil {
							continue
						}

						// this to deal with relative links in the HTML doc
						if len(uChild.Scheme) == 0 {
							// extract domain from parent and prefix with it
							childURL = uThis.Scheme + "://" + uThis.Host + childURL
						} else if uChild.Scheme == "mailto" || uChild.Scheme == "monitor" {
							log.Printf("Skipping unsupported scheme '%s' at page '%s'\n", uChild.Scheme, thisURL)
							continue
						}

						// if this crawl is supposed to stay within the original domain
						// then don't venture outside for any explicitly stated
						// external domains
						if domainsticky {
							startDomain := uThis.Hostname()
							childDomain := uChild.Hostname()
							if len(childDomain) > 0 && startDomain != childDomain {
								continue
							}
						}

						URLs = append(URLs, childURL)
					}
				}
			}
		}
	}

	return URLs
}
