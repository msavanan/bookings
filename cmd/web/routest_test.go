package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/msavanan/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var a *config.AppConfig
	mux := routes(a)

	switch v := mux.(type) {

	case *chi.Mux:

	default:
		t.Errorf(fmt.Sprintf("Type is not *chi.mux, type is %T", v))
	}

}
