package main

import (
	"flag"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// hold list of visited URLs
var known = struct {
	sync.RWMutex
	m map[string]bool
}{m: make(map[string]bool)}

type Result struct {
	err       error
	depth     int
	parentURL *string // optional to deal with starting URL having no parent
	thisURL   string
	response  *http.Response
	timeMS    int64
	start     time.Time
	end       time.Time
}

// after parsing command line kick off potentially recursive search over the URLs

func main() {
	maxdepth := flag.Int("maxdepth", 0, "Maximum depth to crawl")
	domainsticky := flag.Bool("domainsticky", true, "If true, stay within original domain")
	timestamps := flag.Bool("timestamps", false, "show start/end timestamps when true")
	flag.Parse()
	URLs := flag.Args()
	cosmeticPlural := ""
	if len(URLs) != 1 {
		cosmeticPlural = "s"
	}
	log.Printf("Starting scan of %d URL%s...\n", len(URLs), cosmeticPlural)

	// open a waitgroup to make sure we close resources only when all workers are done
	// use two channels:  reporting for workers to share their results;  command to talk
	// to the collector

	var wg sync.WaitGroup
	reporting := make(chan Result, 10)
	command := make(chan string)
	go Collector(reporting, command, *timestamps)

	start := time.Now()
	for _, url := range URLs {
		wg.Add(1)
		go visit(reporting, &wg, nil, url, 1, *maxdepth, *domainsticky)
	}

	// block for all URLs to be visited
	wg.Wait()
	duration := time.Now().Sub(start)
	inMillis := duration.Nanoseconds() / int64(time.Millisecond)

	// tell collator to clean up and wait here for the response that it did
	command <- "quit"
	<-command

	close(reporting)
	close(command)

	log.Printf("Found %d unique URLs in %d ms.\n", len(known.m), inMillis)
}

func visit(reporting chan Result, wg *sync.WaitGroup, parentURL *string, thisURL string, depth int, maxdepth int, domainsticky bool) {

	defer wg.Done()

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

	start := time.Now()
	response, err := http.Get(thisURL)
	end := time.Now()
	duration := end.Sub(start)
	inMillis := duration.Nanoseconds() / int64(time.Millisecond)
	result := Result{err, depth, parentURL, thisURL, response, inMillis, start, end}
	reporting <- result

	if err == nil {
		// maxdepth of zero means: no max depth , i.e. infinite
		if maxdepth != 0 && depth >= maxdepth {
			return
		}

		// don't continue if this looks like application data such as a PDF
		header := response.Header
		if header["Content-Type"] != nil {
			for _, val := range header["Content-Type"] {
				if strings.Contains(val, "application") {
					log.Printf("skipping Content-Type '%s' for page '%s'\n", val, thisURL)
					return
				}
			}
		}

		body := response.Body
		defer body.Close()
		URLs := extractHyperlinks(thisURL, html.NewTokenizer(body), domainsticky)
		for _, url := range URLs {
			wg.Add(1)
			go visit(reporting, wg, &thisURL, url, depth+1, maxdepth, domainsticky)
		}
	}
}
