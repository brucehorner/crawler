package main

import "testing"

// kind of risky testing channels like this in case the goroutine actually
// does have a defect and blocks the test execution

func WillItObeyCommands(t *testing.T) {
	command := make(chan string)
	dummy := make(chan Result)
	go Collector(dummy, command, false)
	command <- "quit"
	res := <-command
	if res != "done" {
		t.Errorf("Expected 'done' from collector and got '%s'", res)
	}
}
