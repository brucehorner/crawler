package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var err = errors.New("oops")
var depth = 2
var parentURL = "parent"
var thisURL = "this"
var status = "200 OK"
var protocol = "HTTP/1.1"
var response = http.Response{status, 200, protocol, 1, 0, nil, nil, 0, nil, false, false, nil, nil, nil}
var timeMS = int64(4235)
var start = time.Now()
var end = time.Now()
var timeStr = strconv.FormatInt(timeMS, 10)
var startStr = start.Format(time.StampMilli)
var endStr = end.Format(time.StampMilli)

func EmptyResult(t *testing.T) {

	var r Result
	res := format(r, false)
	target := ""
	if res != target {
		t.Errorf("Expected empty string for nil Result, got '%s'", res)
	}
}

func ErrorCondition(t *testing.T) {

	res := format(Result{err, depth, &parentURL, thisURL, &response, timeMS, start, end}, false)
	target := fmt.Sprintf("%2d ", depth) + parentURL + " " + thisURL + " \"" + err.Error() + "\" " + timeStr
	if res != target {
		t.Errorf("Expected '%s' but got '%s'", target, res)
	}
}

func ErrorConditionWithTimestamps(t *testing.T) {

	res := format(Result{err, depth, &parentURL, thisURL, &response, timeMS, start, end}, true)
	target := fmt.Sprintf("%2d ", depth) + parentURL + " " + thisURL + " \"" + err.Error() + "\" " + timeStr + " \"" + startStr + "\" \"" + endStr + "\""
	if res != target {
		t.Errorf("Expected '%s' but got '%s'", target, res)
	}
}

func TestNilParent(t *testing.T) {

	res := format(Result{err, depth, nil, thisURL, &response, timeMS, start, end}, false)
	target := fmt.Sprintf("%2d ", depth) + "root " + thisURL + " \"" + err.Error() + "\" " + timeStr
	if res != target {
		t.Errorf("Expected '%s' but got '%s'", target, res)
	}
}

func TestRegularCondition(t *testing.T) {

	res := format(Result{nil, depth, &parentURL, thisURL, &response, timeMS, start, end}, false)
	target := fmt.Sprintf("%2d ", depth) + parentURL + " " + thisURL + " \"" + status + "\" " + timeStr
	if res != target {
		t.Errorf("Expected '%s' but got '%s'", target, res)
	}
}
