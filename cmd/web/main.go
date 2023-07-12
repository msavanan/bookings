package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/driver"
	"github.com/msavanan/bookings/internal/handlers"
	"github.com/msavanan/bookings/internal/helpers"
	"github.com/msavanan/bookings/internal/models"
	"github.com/msavanan/bookings/internal/render"
)

const (
	port = ":8080"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

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

func run() (*driver.DB, error) {

	//What am i going to store other than basic types in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	//gob.Register(models.RoomRestriction{})
	gob.Register(models.Room{})

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("Connecting to database....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=12345678")
	if err != nil {
		log.Fatal("Can't connect to database")
		return nil, err
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("Cannot create template cache", err)
		return db, err
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction //false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelper(&app)

	return db, nil
}
