package main

import (
	"fmt"
	"strings"
	"time"
)

// format results structure into a 'pretty' string for printing and loading
// as space delimited file for analyss
//
// 1 - depth
// 2 - parent URL
// 3 - this URL
// 4 - query status string OR error message if failure
// 5 - response time in ms
// 6 - start wall clock time of the query, in ms
// 7 - end wall clock time, in ms
//
func format(result Result, timestamps bool) string {

	var strs strings.Builder
	strs.WriteString(fmt.Sprintf("%2d ", result.depth))
	if result.parentURL == nil {
		strs.WriteString("root")
	} else {
		strs.WriteString(*result.parentURL)
	}
	strs.WriteString(fmt.Sprintf(" %s", result.thisURL))

	var response string
	if result.err == nil && result.response != nil {
		response = result.response.Status
	} else {
		response = result.err.Error()
	}
	strs.WriteString(fmt.Sprintf(" \"%s\" %d", response, result.timeMS))

	if timestamps {
		strs.WriteString(fmt.Sprintf(" \"%s\" \"%s\"", result.start.Format(time.StampMilli), result.end.Format(time.StampMilli)))
	}

	return strs.String()
}
