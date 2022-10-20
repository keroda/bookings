package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keroda/bookings/internal/config"
	"github.com/keroda/bookings/internal/driver"
	"github.com/keroda/bookings/internal/handlers"
	"github.com/keroda/bookings/internal/helpers"
	"github.com/keroda/bookings/internal/models"
	"github.com/keroda/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	defer close(app.MailChan)
	listenForMail()
	fmt.Println("Started mail listener")

	fmt.Printf("Started app on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*pgxpool.Pool, error) {
	//what to put into the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	//read flags
	// inProduction := flag.Bool("production", true, "Application is in production")
	// useCache := flag.Bool("cache", true, "Use template cache")
	// dbHost := flag.String("dbhost", "localhost", "Database host")
	// dbName := flag.String("dbname", "", "Database name")
	// dbUser := flag.String("dbuser", "", "Database user")
	// dbPass := flag.String("dbpass", "", "Database password")
	// dbPort := flag.String("dbport", "5432", "Database port")
	// dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")
	// flag.Parse()
	// connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	// if *dbName == "" {
	// 	fmt.Println("Missing required flags")
	// 	os.Exit(1)
	// }

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//app.InProduction = *inProduction //from flag
	//app.UseCache = *useCache

	app.InProduction = false //*inProduction //from flag
	app.UseCache = false     //*useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("Connecting to database...")

	//db, err := driver.ConnectSQL(connectionString)  //<<<<----------------------
	db, err := driver.ConnectSQL(driver.MyDb)
	if err != nil {
		log.Fatal("Cannot connect to database!", err)
	}
	fmt.Println("Connected to database")

	tc, err := render.CreateTemplatecache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
