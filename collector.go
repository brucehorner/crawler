package main

import (
	"fmt"
	"log"
)

// this is a simple console printing worker result collector
// it listens on the reporting channel for interesting things from
// the workers, and on the command channel for instructions which
// for now is only when to clean up and shutdown
func Collector(reporting chan Result, command chan string, timestamps bool) {

	for {
		select {

		case message := <-reporting:
			prettyString := format(message, timestamps)
			fmt.Println(prettyString)

		case <-command:
			// there is no actual clean up for this collector
			log.Println("Collector is done.")
			command <- "done"
		}
	}

}
