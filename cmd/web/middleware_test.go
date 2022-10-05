package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		//it's good , so do nothing
	default:
		t.Errorf("type is not http.handler, it's a %T", v)

	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		//it's good , so do nothing
	default:
		t.Errorf("type is not http.handler, it's a %T", v)

	}
}
