package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/handlers"
	"github.com/msavanan/bookings/internal/models"
	"github.com/msavanan/bookings/internal/render"
)

const (
	port = ":8080"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server at port:", port)
	//http.ListenAndServe(port, nil)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	//What am i going to store other than basic types in session
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("Cannot create template cache", err)
		return err
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction //false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	return nil
}
