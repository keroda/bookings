package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/keroda/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//all good
	default:
		t.Errorf("expected chi.Mux but got %T", v)
	}
}
