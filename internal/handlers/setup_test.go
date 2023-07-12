package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/models"
	"github.com/msavanan/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

const pathToTemplates = "./../../templates"

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Println("Cannot create template cache", err)
	}

	app.TemplateCache = tc
	app.UseCache = true //app.InProduction //false

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf) //CSRF Token
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)

	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/generals-quarters", Repo.Generals)

	//GET service Availablity
	mux.Get("/search-availability", Repo.Availability)

	//POST service Availablity
	mux.Post("/search-availability", Repo.PostAvailability)

	//json service Availablity
	mux.Post("/search-availability-json", Repo.AvailabilityJson)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	//get all of the files named *,page.html from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	// range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts

	}
	return cache, nil

}
