package main

import (
	"flag"
	"time"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
  "sync"
)

var known = struct{
  sync.RWMutex
  m map[string]bool
}{m: make(map[string]bool)}

func main() {
	maxdepth := flag.Int("maxdepth", 0, "Maximum depth to crawl")
  domainsticky := flag.Bool("domainsticky", true, "If true, stay within original domain")
	flag.Parse()
	URLs := flag.Args()
  var waitgroup sync.WaitGroup
//TODO:  finish concurrency
//  waitgroup.Add(len(URLs))
	for _, url := range URLs {
		process(&waitgroup, nil, url, 1, *maxdepth, *domainsticky)
		fmt.Println()
	}
//  waitgroup.Wait()
  fmt.Printf("Completed. Found %d unique URLs\n", len(known.m))
}


func process(waitgroup *sync.WaitGroup, parentURL *string, thisURL string, depth int, maxdepth int, domainsticky bool) {

//  defer waitgroup.Done()

  // bail out if this has already been seen
  // otherwise add this to the found list
  known.RLock()
  if known.m[thisURL] {
    known.RUnlock()
    return
	} else {
    known.RUnlock()
    known.Lock()
		known.m[thisURL] = true
    known.Unlock()
	}

	fmt.Printf("%3d ", depth)
	if parentURL == nil {
		fmt.Print("root")
	} else {
		fmt.Print(*parentURL)
	}
	fmt.Printf(" %s", thisURL)

  start := time.Now()
	response, error := http.Get(thisURL)
  duration := time.Now().Sub(start)
	if error == nil {
    inMillis := duration.Nanoseconds() / int64(time.Millisecond)
		fmt.Printf(" \"%s\" %dms\n", response.Status, inMillis)

	// maxdepth of zero means: no max depth , i.e. infinite
		if maxdepth != 0 && depth >= maxdepth {
			return
		}

		body := response.Body
		defer body.Close()
		tokenizer := html.NewTokenizer(body)
		for {
			tokenType := tokenizer.Next()
			switch {
			case tokenType == html.ErrorToken:
				return
			case tokenType == html.StartTagToken:
				token := tokenizer.Token()
        // skip every tag except anchor tags
				if token.Data != "a" {
					continue
				} else {
					for _, a := range token.Attr {
						if a.Key == "href" {
							childURL := a.Val
							if len(childURL) == 0 {
								continue
							}

							uThis, _ := url.Parse(thisURL)
							uChild, e := url.Parse(childURL)
              if e != nil {
								fmt.Println(childURL, " =>", e)
								continue
							}
							if len(uChild.Scheme) == 0 {
								// extract domain from parent and prefix with it
								childURL = uThis.Scheme + "://" + uThis.Host + childURL
							} else if uChild.Scheme == "mailto" {
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

							//waitgroup.Add(1)
							process(waitgroup, &thisURL, childURL, depth + 1, maxdepth, domainsticky)
						}
					}
				}
			}
		}
	} else {
		fmt.Printf(" %s\n", error)
	}
}
