package main

import "testing"

func TestNilParameters(t *testing.T) {

	res := extractHyperlinks("", nil, false)
	if res != nil || len(res) != 0 {
		t.Errorf("Expected empty array for nil body, but got '%s'", res)
	}
}
