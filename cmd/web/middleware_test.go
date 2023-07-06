package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	mH := myHandler{}

	h := NoSurf(&mH)

	switch v := h.(type) {
	case http.Handler:

	default:
		t.Errorf(fmt.Sprintf("Error in handler, type is NOT http.Handler but the type is %T", v))
	}

}

func TestSessionLoad(t *testing.T) {

	mH := myHandler{}

	h := SessionLoad(&mH)

	switch v := h.(type) {
	case http.Handler:

	default:
		t.Errorf(fmt.Sprintf("Error in handler, type is NOT http.Handler but the type is %T", v))
	}

}
