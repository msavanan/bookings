package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/msavanan/bookings/pkg/config"
	"github.com/msavanan/bookings/pkg/handlers"
	"github.com/msavanan/bookings/pkg/render"
)

const (
	port = ":8080"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	//session.Cookie.Persist = true
	//session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction //false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	log.Println("Starting server at port:", port)
	//http.ListenAndServe(port, nil)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	srv.ListenAndServe()
}
