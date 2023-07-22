package main

import (
	"encoding/gob"
	"flag"
	"fmt"
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
	defer close(app.MailChan)

	log.Println("Listening for email......")
	listenForMail()

	// log.Println("starting to send email.......")
	// from := "me@here.com"
	// auth := smtp.PlainAuth("", from, "", "localhost")
	// err = smtp.SendMail("localhost:1025", auth, from, []string{"yo@there.com"}, []byte("Hello World!"))
	// if err != nil {
	// 	log.Println("failed to sent email via smtp package ", err)
	// } else {
	// 	log.Println("Sent email successfully...")
	// }

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
	gob.Register(map[string]int{})

	//Flags
	inProduction := flag.Bool("production", true, "Application in production")
	useCache := flag.Bool("cache", true, "use template cache")
	dbHost := flag.String("dbhost", "", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings(disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailchan := make(chan models.MailData)

	app.MailChan = mailchan

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InProduction = *inProduction

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("Connecting to database....")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
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
	app.UseCache = *useCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelper(&app)

	return db, nil
}
