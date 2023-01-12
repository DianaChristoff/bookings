package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DianaChristoff/bookings/pkg/config"
	"github.com/DianaChristoff/bookings/pkg/handlers"
	"github.com/DianaChristoff/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"

var app config.AppConfig
var session *scs.SessionManager

// main application function
func main() {

	app.InProduction = false

	//Creating new session that lasts 24 hours and persists(is not destroyed when closing the browser)
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	//how strict do you want to be when applying cookie to site
	session.Cookie.SameSite = http.SameSiteLaxMode
	//using port 3000 so it's not secure
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
