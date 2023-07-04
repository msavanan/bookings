package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/handlers"
)

func routes(a *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf) //CSRF Token
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/generals-quarters", handlers.Repo.Generals)

	//GET service Availablity
	mux.Get("/search-availability", handlers.Repo.Availability)

	//POST service Availablity
	mux.Post("/search-availability", handlers.Repo.PostAvailability)

	//json service Availablity
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJson)

	mux.Get("/make-reservation", handlers.Repo.Reservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
